package redisutil

import (
	"context"
	"encoding/json"
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
		return err
	}

	return json.Unmarshal([]byte(value), marshalTo)
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
				logrus.Error("failed to unmarshal Game State from key", err)
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
				logrus.Error("failed to unmarshal Game State from pubsub", err)
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
	client *redis.Client, key string, publishKey string, ttl time.Duration, maxRetryDuration time.Duration,
	defaultValuePtr interface{},
	update func(value interface{}) interface{}) (err error) {
	operation := func() error {
		return client.Watch(ctx, func(tx *redis.Tx) error {
			var valueJSON string

			valueJSON, err = tx.Get(ctx, key).Result()

			if err != nil {
				if err != redis.Nil {
					return backoff.Permanent(fmt.Errorf("failed to get current value from redis during update transaction %w", err))
				}
			} else {
				err = json.Unmarshal([]byte(valueJSON), defaultValuePtr)
				if err != nil {
					return backoff.Permanent(fmt.Errorf("failed to deserialize current value during update transaction %w", err))
				}
			}

			newValue := update(defaultValuePtr)
			newValueSerialized, err := json.Marshal(newValue)
			if err != nil {
				return backoff.Permanent(fmt.Errorf("failed to serialize new value during update transaction %w", err))
			}
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, newValueSerialized, ttl)
				if len(publishKey) > 0 {
					pipe.Publish(ctx, publishKey, newValueSerialized)
				}
				return nil
			})
			if err != nil {
				if err == redis.TxFailedErr {
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
