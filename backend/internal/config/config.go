package config

import (
	"os"
	"sync"
)

// Config 單例 env 配置
type ConfigStruct struct {
	Port           string
	RedisURL       string
	KafkaURL       string
	SqlUrl         string
	FrontendOrigin string
	ReserveTopic   string
	TempoURL       string
}

var (
	config     *ConfigStruct
	configonce sync.Once
)

// GetInstance 取得單例
func GetConfig() *ConfigStruct {
	configonce.Do(func() {
		c := &ConfigStruct{
			ReserveTopic:   getEnv("RESERVE_TOPIC", "reserve"),
			FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
			Port:           getEnv("PORT", "8080"),
			RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
			KafkaURL:       getEnv("KAFKA_URL", "localhost:9092"),
			TempoURL:       getEnv("TEMPO_URL", "localhost:4318"),
			SqlUrl:         getEnv("SQL_URL", "postgres://postgres:postgres@localhost:5432/testdb"),
		}
		config = c
	})
	return config
}

// getEnv 若環境變數不存在，回傳預設值
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
