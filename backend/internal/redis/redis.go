package redis

import (
	"backend/internal/config"
	"backend/internal/structs"
	"context"
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

const seat_set = "seat_set"

func InitRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.GetConfig().RedisURL,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect or ping Redis: %w", err)
	}
	fmt.Printf("Redis connected successfully: %s\n", pong)
	RedisClient = rdb
	loadScripts()
	return nil
}
func createRedisKey(seatId int) string {
	return fmt.Sprintf("{seat}:%d", seatId)
}
func getSeatLockKeyList(reserveData *structs.ReservePostData) []string {
	var lockKeys []string
	for _, seatId := range reserveData.SeatIds {
		// 使用 fmt.Sprintf 進行格式化
		lockKey := createRedisKey(seatId)
		lockKeys = append(lockKeys, lockKey)
	}
	return lockKeys
}
func GetAllLockedSeatIDs() ([]int, error) {
	// 1️⃣ 取得集合所有元素
	var ctx = context.Background()
	keys, err := RedisClient.SMembers(ctx, seat_set).Result()
	if err != nil {
		return nil, err
	}

	var seatIDs []int
	for _, key := range keys {
		// 2️⃣ key 是 "{seat}:123"，我們要取得 123
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue // 異常格式就跳過
		}

		id, err := strconv.Atoi(parts[1])
		if err != nil {
			continue // 轉整數失敗就跳過
		}

		seatIDs = append(seatIDs, id)
	}

	return seatIDs, nil
}

func CleanExpire() {
	log.Println("✅ 開始清理過期座位...")
	ctx := context.Background()
	totalRemoved := 0

	// 1️⃣ 一次性取得 seat_set 所有元素
	keys, err := RedisClient.SMembers(ctx, seat_set).Result()
	if err != nil {
		log.Printf("❌ CleanExpire 取得 seat_set 失敗: %v", err)
		return
	}

	for _, key := range keys {
		// 2️⃣ 判斷鍵是否已過期
		exists, err := RedisClient.Exists(ctx, key).Result()
		if err != nil {
			log.Printf("⚠️ EXISTS 判斷錯誤, key: %s, err: %v", key, err)
			continue
		}

		if exists == 0 {
			// 3️⃣ 如果鍵不存在，從 seat_set 移除
			if _, err := RedisClient.SRem(ctx, seat_set, key).Result(); err != nil {
				log.Printf("⚠️ SREM 移除鍵失敗, key: %s, err: %v", key, err)
			} else {
				totalRemoved++
				log.Printf("✅ 已移除過期鍵: %s", key)
			}
		}
	}

	log.Printf("🔹 CleanExpire 完成, 共移除 %d 個過期鍵", totalRemoved)
}
