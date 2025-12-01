package kafka

import (
	config "backend/internal/config"
	localOtel "backend/internal/otel"
	"backend/internal/redis"
	"backend/internal/sql"
	"backend/internal/structs"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func InitKafkaReader(ctx context.Context, wg *sync.WaitGroup) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.GetConfig().KafkaURL},
		GroupID: "consumer",
		Topic:   config.GetConfig().ReserveTopic,
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := r.Close(); err != nil {
				log.Println("failed to close reader:", err)
			}
		}()

		for {
			m, err := r.ReadMessage(ctx) // 用傳入的 context
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Kafka reader exiting due to context cancel")
					// execute defer
					return
				}
				log.Println("ReadMessage error:", err)
				continue
			}
			handleConsumeKafka(&m)
		}
	}()
}
func handleConsumeKafka(m *kafka.Message) {
	// 1. Extract context from Kafka headers
	carrier := propagation.MapCarrier{}
	for _, h := range m.Headers {
		carrier[h.Key] = string(h.Value)
	}
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

	// 2. Consumer root span
	ctx, span := localOtel.GlobalTracer.Start(ctx, "ConsumeReserveMessage")
	defer span.End()

	var reserveData structs.ReservePostData
	_, unmarshalSpan := localOtel.GlobalTracer.Start(ctx, "JSON_Unmarshal")
	if err := json.Unmarshal(m.Value, &reserveData); err != nil {
		unmarshalSpan.RecordError(err)
		unmarshalSpan.End()
		log.Printf("❌ JSON 解析錯誤: %v, 原始內容: %s", err, string(m.Value))
		return
	}
	unmarshalSpan.End()

	// 3. SQL 更新
	_, sqlSpan := localOtel.GlobalTracer.Start(ctx, "SQL_ReserveSeats")
	if err := sql.ReserveSeats(&reserveData); err != nil {
		sqlSpan.RecordError(err)
		sqlSpan.End()
		log.Printf("❌ SQL 更新失敗: %v", err)
		if err2 := redis.ReleaseLockSeats(&reserveData); err2 != nil {
			log.Printf("⚠️ 釋放 Redis 鎖失敗: %v", err2)
		}
		return
	}
	sqlSpan.End()

	// 4. Redis 更新
	_, redisSpan := localOtel.GlobalTracer.Start(ctx, "Redis_HandleReserve")
	if err := redis.HandleReserve(&reserveData); err != nil {
		redisSpan.RecordError(err)
		redisSpan.End()

		// 回滾
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
	redisSpan.End()

	log.Printf("✅ 消費完成並更新 SQL+Redis 成功, ID: %d", reserveData.UserId)
}
