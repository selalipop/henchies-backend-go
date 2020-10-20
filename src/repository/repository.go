package repository

import "github.com/go-redis/redis/v8"

// Repository allows for CRUD operations on models
type Repository struct {
	RedisClient *redis.Client
}
