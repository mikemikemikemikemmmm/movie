package main

import (
	"backend/internal"
	"backend/internal/kafka"
	"backend/internal/redis"
	"backend/internal/sql"
	"log"
	"time"
)

func main() {
	if err := kafka.InitKafkaWriter(); err != nil {
		log.Printf("kafka初始化失敗 : %v", err)
	}
	go kafka.InitKafkaReader()
	if err := redis.InitRedis(); err != nil {
		log.Printf("redis初始化失敗 : %v", err)
	}
	if err := sql.InitSQL(); err != nil {
		log.Printf("Database初始化失敗 : %v", err)
	}
	internal.InitRouter()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		redis.CleanExpire()
	}
}
