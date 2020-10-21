package redisutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

// GetRedisJSON retrieves a JSON serialized value from Redis, marshalTo must be a pointer
func GetRedisJSON(ctx context.Context, client *redis.Client, key string, marshalTo interface{}) error {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get json redis key (%v): %w", key, err)
	}

	if err = json.Unmarshal([]byte(value), marshalTo); err != nil {
		return fmt.Errorf("failed to unmarshal json retrieved from redis key (%v): %w", key, err)
	}
	return nil
}

// SubscribeJSON retrieves a JSON serialized value from Redis and listens for updates to it's value, marshalTo must be a pointer
// The initial value is fetched and sent using getKey, use an empty string to disable this behavior
func SubscribeJSON(ctx context.Context, client *redis.Client, getKey string, pubSubKey string, marshalTo interface{}) (channel chan interface{}, err error) {
	listen := client.Subscribe(ctx, pubSubKey).Channel()

	send := make(chan interface{})

	go func() {
		defer close(send)
		if getKey != "" {
			err := GetRedisJSON(ctx, client, getKey, marshalTo)
			if err != nil {
				logrus.Error("failed to unmarshal game state from key", err)
				return
			}
			send <- marshalTo
		}

		for {
			message, ok := <-listen
			if !ok {
				return
			}

			err := json.Unmarshal([]byte(message.Payload), marshalTo)
			if err != nil {
				logrus.Error("failed to unmarshal game state from pubsub", err)
				break
			}
			send <- marshalTo
		}
	}()

	return send, nil
}

// UpdateKeyTransaction updates a JSON serialized value from Redis transactionally
// The Update Method will be called multiple times if the value is modified by another process
func UpdateKeyTransaction(ctx context.Context,
	client *redis.Client, key string, publishKey string, ttl time.Duration,
	maxRetryDuration time.Duration, defaultValuePtr interface{}, update func(value interface{}) interface{}) (err error) {
	operation := func() error {
		return client.Watch(ctx, func(tx *redis.Tx) error {
			err = GetRedisJSON(ctx, client, key, defaultValuePtr)

			if err != nil {
				if !errors.Is(err, redis.Nil) {
					return backoff.Permanent(fmt.Errorf("failed to get current value from redis during update transaction: %w", err))
				}
			}

			newValue := update(defaultValuePtr)

			_, err = tx.TxPipelined(ctx, getPiperlinerForValue(ctx, newValue, key, publishKey, ttl))

			if err != nil {
				if errors.Is(err, redis.TxFailedErr) {
					return fmt.Errorf("failed to to run update transaction pipeline due to value change %w", err)
				}
				return backoff.Permanent(fmt.Errorf("failed to to run update transaction pipeline %w", err))
			}
			return nil
		}, key)
	}

	txBackoff := backoff.NewExponentialBackOff()
	txBackoff.MaxElapsedTime = maxRetryDuration
	return backoff.Retry(operation, txBackoff)
}

func getPiperlinerForValue(ctx context.Context, newValue interface{}, key string, publishKey string, ttl time.Duration) func(pipe redis.Pipeliner) error {
	return func(pipe redis.Pipeliner) error {
		if newValue == nil {
			pipe.Del(ctx, key)
			if len(publishKey) > 0 {
				pipe.Del(ctx, publishKey)
			}
			return nil
		}
		newValueSerialized, err := json.Marshal(newValue)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to serialize new value during update transaction %w", err))
		}

		pipe.Set(ctx, key, newValueSerialized, ttl)
		if len(publishKey) > 0 {
			pipe.Publish(ctx, publishKey, newValueSerialized)
		}
		return nil
	}
}
