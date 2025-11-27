package redis

import (
	"backend/internal/structs"
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ReleaseLockSeats(reserveData *structs.ReservePostData) error {

	lockKeys := getSeatLockKeyList(reserveData)
	lockVal := strconv.Itoa(reserveData.UserId)
	_, err := releaseLockScript.Run(
		context.Background(),
		RedisClient,
		lockKeys,
		lockVal,
	).Result()

	if err != nil && err != redis.Nil {
		// 處理執行錯誤 (例如連線問題)
		return fmt.Errorf("failed to release lock script: %w", err)
	}
	return nil
}
