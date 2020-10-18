package repository

import "github.com/go-redis/redis/v8"
import "context"

type RepositoryEnv struct {
	RedisClient *redis.Client
	Context     context.Context
}
