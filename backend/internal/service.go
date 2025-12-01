package internal

import (
	"backend/internal/kafka"
	"backend/internal/otel"
	"backend/internal/redis"
	"backend/internal/sql"
	"backend/internal/structs"
	"log"

	"github.com/gin-gonic/gin"
)

func handleReserve(c *gin.Context) {
	c.Header(otel.NotEndSpan, "true")
	ctx := c.Request.Context() // 從 middleware 傳下來的 context

	// 解析 JSON
	var reserveData structs.ReservePostData
	_, span := otel.GlobalTracer.Start(ctx, "BindJSON")
	if err := c.ShouldBindJSON(&reserveData); err != nil {
		span.RecordError(err)
		span.End()
		c.JSON(400, gin.H{"error": "請求參數格式錯誤: " + err.Error()})
		return
	}
	span.End()

	// 驗證 SeatIds
	if len(reserveData.SeatIds) == 0 {
		c.JSON(400, gin.H{"error": "SeatIds 不能為空"})
		return
	}

	// Redis 鎖座位
	_, span = otel.GlobalTracer.Start(ctx, "RedisLockSeats")
	if err := redis.HandleLockSeats(&reserveData); err != nil {
		span.RecordError(err)
		span.End()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	span.End()

	// Kafka 發送訂位消息
	_, span = otel.GlobalTracer.Start(ctx, "KafkaSendReserveMessage")
	if err := kafka.SendReserveMessage(c, &reserveData); err != nil {
		span.RecordError(err)
		span.End()
		redis.ReleaseLockSeats(&reserveData)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	span.End()

	c.JSON(200, gin.H{"status": "success", "message": "訂位請求已排隊"})
}
func handleGetAllSeats(c *gin.Context) {
	seats, err := sql.GetAllSeats()
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"error": err})
		return
	}
	for i := range seats {
		seats[i].UserID = nil
	}
	lockedSeatIDList, err := redis.GetAllLockedSeatIDs()
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"error": err})
		return
	}
	lockedMap := make(map[uint]struct{}, len(lockedSeatIDList))
	for _, id := range lockedSeatIDList {
		lockedMap[uint(id)] = struct{}{} // 將 int 轉成 uint
	}

	for i := range seats {
		if _, ok := lockedMap[seats[i].ID]; ok {
			seats[i].Status = "reserved"
		}
	}
	c.JSON(200, gin.H{"data": seats})
}
func handleRefreshSeats(c *gin.Context) {
	lockedSeatIDList, err := redis.GetAllLockedSeatIDs()
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"data": lockedSeatIDList})
}
func handleCheckReserve(c *gin.Context) {
	var reserveData structs.ReservePostData
	if err := c.ShouldBindJSON(&reserveData); err != nil {
		c.JSON(400, gin.H{
			"error": "請求參數格式錯誤: " + err.Error(),
		})
		return
	}
	result := redis.CheckReserve(&reserveData)
	switch result {
	case "success":
		c.JSON(200, gin.H{"data": "訂位成功"})
	case "processing":
		c.JSON(202, gin.H{"data": "訂位處理中"})
	default:
		c.JSON(400, gin.H{"data": "訂位失敗"})
	}
}
func handleCheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{"data": "health"})
}

func handleCheckReady(c *gin.Context) {
	// 檢查 SQL 連線
	if err := sql.CheckSqlReady(); err != nil {
		c.JSON(500, gin.H{"error": "sql not ready", "detail": err.Error()})
		return
	}

	// 檢查 Redis
	if err := redis.RedisClient.Ping(c).Err(); err != nil {
		c.JSON(500, gin.H{"error": "redis not ready", "detail": err.Error()})
		return
	}

	// 檢查 Kafka
	if err := kafka.CheckKafkaReady(); err != nil {
		c.JSON(500, gin.H{"error": "kafka not ready", "detail": err.Error()})
		return
	}

	// 若全部健康
	c.JSON(200, gin.H{"data": "ready"})
}
