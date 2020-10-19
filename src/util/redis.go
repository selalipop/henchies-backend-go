package util

import (
	"context"
	"encoding/json"
	. "github.com/go-redis/redis/v8"
)

func GetRedisJson(ctx context.Context, client *Client, key string, marshalTo interface{}) error {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), marshalTo)
}
