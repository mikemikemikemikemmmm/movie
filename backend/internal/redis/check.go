package redis

import (
	"backend/internal/structs"
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func CheckReserve(reserveData *structs.ReservePostData) string {
	lockKeys := getSeatLockKeyList(reserveData)
	lockVal := strconv.Itoa(reserveData.UserId)
	result, err := checkScript.Run(
		context.Background(),
		RedisClient,
		lockKeys, // KEYS
		lockVal,  // ARGV[1]
	).Result()

	if err != nil && err != redis.Nil {
		log.Printf("check reserve redis連線錯誤: %v", err)
		return "error"
	}

	switch val := result.(type) {
	case string:
		return val
	default:
		log.Printf("check reserve redis 非string錯誤: %v", err)
		return "error"
	}
}
