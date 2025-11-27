package redis

import (
	"backend/internal/structs"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func HandleLockSeats(reserveData *structs.ReservePostData) error {
	fmt.Println("開始鎖seats")
	lockKeys := getSeatLockKeyList(reserveData)
	lockVal := strconv.Itoa(reserveData.UserId)
	result, err := lockSeatScript.Run(
		context.Background(),
		RedisClient,
		lockKeys,
		lockVal,
	).Result()

	if err != nil && err != redis.Nil {
		// 處理執行錯誤 (例如連線問題)
		log.Printf("failed to run multi-lock script: %v", err)
		return fmt.Errorf("資料庫連線錯誤")
	}
	switch val := result.(type) {
	case string:
		if val == "success" {
			fmt.Println("鎖seats成功")
			return nil
		} else {
			return fmt.Errorf("部分位子已被訂走")
		}
	default:
		// 處理所有其他非預期的回傳類型
		log.Printf("failed to run multi-lock script: %v", err)
		return fmt.Errorf("訂位時發生錯誤")
	}
}
