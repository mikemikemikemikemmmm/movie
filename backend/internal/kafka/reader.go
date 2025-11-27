package kafka

import (
	config "backend/internal/config"
	"backend/internal/redis"
	"backend/internal/sql"
	"backend/internal/structs"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func InitKafkaReader() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.GetConfig().KafkaURL},
		GroupID: "consumer",
		Topic:   config.GetConfig().ReserveTopic,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		handleConsumeKafka(&m)
	}

	if err := r.Close(); err != nil {
		log.Println("failed to close reader:", err)
	}
}
func handleConsumeKafka(m *kafka.Message) {
	var reserveData structs.ReservePostData
	if err := json.Unmarshal(m.Value, &reserveData); err != nil {
		log.Printf("❌ JSON 解析錯誤: %v, 原始內容: %s", err, string(m.Value))
		return
	}

	// 1. 嘗試更新 SQL
	if err := sql.ReserveSeats(&reserveData); err != nil {
		log.Printf("❌ SQL 更新失敗: %v", err)
		if err2 := redis.ReleaseLockSeats(&reserveData); err2 != nil {
			log.Printf("⚠️ 釋放 Redis 鎖失敗: %v", err2)
		}
		return
	}

	// 2. 嘗試更新 Redis
	if err := redis.HandleReserve(&reserveData); err != nil {
		log.Printf("❌ Redis 更新失敗: %v", err)

		errSQL := sql.RollbackReserveSeats(&reserveData)
		errRedis := redis.ReleaseLockSeats(&reserveData)

		switch {
		case errSQL != nil && errRedis != nil:
			log.Printf("⚠️ SQL 回滾 & Redis 回滾都失敗，請手動處理")
		case errSQL != nil:
			log.Printf("⚠️ SQL 回滾失敗，請手動處理")
		case errRedis != nil:
			log.Printf("⚠️ Redis 回滾失敗，請手動處理")
		default:
			log.Printf("✅ Redis 更新失敗，回滾成功")
		}
		return
	}

	log.Printf("✅ 消費完成並更新 SQL+Redis 成功, ID: %d", reserveData.UserId)
}
