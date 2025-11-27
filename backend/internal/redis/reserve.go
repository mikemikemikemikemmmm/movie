package redis

import (
	"backend/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type ScriptResult struct {
	Status string `json:"status"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func HandleReserve(reserveData *structs.ReservePostData) error {
	ctx := context.Background()

	keys := getSeatLockKeyList(reserveData)
	val := strconv.Itoa(reserveData.UserId)

	raw, err := reserveScript.Run(ctx, RedisClient, keys, val).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("HandleReserve: redis script execution failed (keys=%v, value=%s): %w",
			keys, val, err)
	}
	fmt.Println(raw, 123)
	// 回傳必須是字串
	rawStr, ok := raw.(string)
	if !ok {
		return fmt.Errorf("HandleReserve: script returned non-string type (keys=%v, type=%T, value=%#v)",
			keys, raw, raw)
	}

	var res ScriptResult
	if err := json.Unmarshal([]byte(rawStr), &res); err != nil {
		return fmt.Errorf("HandleReserve: failed to unmarshal script JSON (keys=%v, raw=%s): %w",
			keys, rawStr, err)
	}

	// success
	if res.Status == "success" {
		return nil
	}

	// failed
	return fmt.Errorf(
		"HandleReserve: script failed (key=%s, value=%s, reason=%s)",
		res.Key, res.Value, res.Reason,
	)
}
