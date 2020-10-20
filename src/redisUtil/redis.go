package redisUtil

import (
	. "context"
	"encoding/json"
	. "github.com/cenkalti/backoff"
	. "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

func GetRedisJson(ctx Context, client *Client, key string, marshalTo interface{}) error {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), marshalTo)
}

func SubscribeJson(ctx Context,client *Client,
	getKey string, pubSubKey string, marhsalTo interface{}) (channel chan interface{}, err error) {

	listen := client.Subscribe(ctx, pubSubKey).Channel()

	send := make(chan interface{})

	go func() {
		defer close(send)
		if getKey != "" {
			err := GetRedisJson(ctx, client, getKey, marhsalTo)
			if err != nil {
				logrus.Error("failed to unmarshal Game State from key", err)
				return
			}
			send <- marhsalTo
		}

		for {
			message, ok := <-listen
			if !ok {
				return
			}

			err := json.Unmarshal([]byte(message.Payload), marhsalTo)
			if err != nil {
				logrus.Error("failed to unmarshal Game State from pubsub", err)
				break
			}
			send <- marhsalTo
		}
	}()

	return send, nil
}

func UpdateKeyTransaction(ctx Context,
	client *Client, key string, publishKey string, ttl time.Duration, maxRetryDuration time.Duration,
	unmarshal func(data []byte)(interface{} ,error),
	init func() interface{},
	update func(value interface{}) interface{}) (err error) {

	operation := func() error {
		return client.Watch(ctx, func(tx *Tx) error {
			var currentValue interface{}
			var valueJson string

			valueJson, err = tx.Get(ctx, key).Result()
			switch err {
			case nil:
				currentValue, err = unmarshal([]byte(valueJson))
				if err != nil {
					return err
				}
			case Nil:
				currentValue = init()
			default:
				return err
			}

			newValue := update(currentValue)
			newValueSerialized, err := json.Marshal(newValue)

			_, err = tx.TxPipelined(ctx, func(pipe Pipeliner) error {
				pipe.Set(ctx, key, newValueSerialized, ttl)
				if len(publishKey) > 0 {
					pipe.Publish(ctx, publishKey, newValueSerialized)
				}
				return nil
			})
			return err
		}, key)

	}

	backoff := NewExponentialBackOff()
	backoff.MaxElapsedTime = maxRetryDuration
	return Retry(operation, backoff)
}
