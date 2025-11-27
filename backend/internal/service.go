package internal

import (
	"backend/internal/kafka"
	"backend/internal/redis"
	"backend/internal/sql"
	"backend/internal/structs"
	"log"

	"github.com/gin-gonic/gin"
)

func handleReserve(c *gin.Context) {
	var reserveData structs.ReservePostData
	if err := c.ShouldBindJSON(&reserveData); err != nil {
		c.JSON(400, gin.H{
			"error": "請求參數格式錯誤: " + err.Error(),
		})
		return
	}
	if len(reserveData.SeatIds) == 0 {
		c.JSON(400, gin.H{"error": "SeatIds 不能為空"})
		return
	}
	if err := redis.HandleLockSeats(&reserveData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := kafka.WriteReserveMessages(c, &reserveData); err != nil {
		redis.ReleaseLockSeats(&reserveData)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
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
