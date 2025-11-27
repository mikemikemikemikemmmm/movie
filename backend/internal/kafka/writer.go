package kafka

import (
	config "backend/internal/config"
	"backend/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer
var reserveTopic = config.GetConfig().ReserveTopic

func InitKafkaWriter() error {
	kafkaUrl := config.GetConfig().KafkaURL
	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaUrl},
	})
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaUrl, reserveTopic, 0)
	if err != nil {
		// 如果無法建立基礎連線，立即返回錯誤
		return fmt.Errorf("failed to dial Kafka broker %s: %w", kafkaUrl, err)
	}
	conn.Close() // 檢查成功後關閉連線
	fmt.Println("kafka 建立主題" + reserveTopic + "成功")
	return nil
}

func WriteReserveMessages(c *gin.Context, reserveData *structs.ReservePostData) error {
	jsonValue, err := json.Marshal(reserveData)
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return fmt.Errorf("解析錯誤")
	}
	kafkaMessage := kafka.Message{
		Topic: reserveTopic,
		Value: []byte(jsonValue),
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // 設定 5 秒超時
	defer cancel()

	err = kafkaWriter.WriteMessages(ctx, kafkaMessage)

	if err != nil {
		log.Printf("Kafka WriteMessages 錯誤: %v", err)
		return fmt.Errorf("訊息發送失敗，請稍後重試")
	}
	log.Println("成功發送所有訂位 ID 訊息到 Kafka")
	return nil
}
