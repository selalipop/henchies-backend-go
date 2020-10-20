package repository

import "github.com/go-redis/redis/v8"

type Repository struct {
	RedisClient *redis.Client
}
