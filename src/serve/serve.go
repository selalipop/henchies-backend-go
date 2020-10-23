package main

import (
	"context"
	"github.com/SelaliAdobor/henchies-backend-go/src/controllers"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func main() {
	args := getArgumentsFromEnv()

	redisClient := getRedisClient(args)

	c := controllers.Controllers{
		Repository: repository.Repository{
			RedisClient: redisClient,
		},
	}

	g := gin.New()

	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.TraceLevel)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	g.Use(ginlogrus.Logger(logger), gin.Recovery(), cors.New(config))

	setupRoutes(g, c)

	if err := g.Run(); err != nil {
		logrus.Panicf("failed to start gin %v", err)
	}
}

func getRedisClient(args Arguments) *redis.Client {
	redisOptions, err := redis.ParseURL(args.RedisConnectURL)
	if err != nil {
		logrus.Panicf("failed to parse redis connection url %v", err)
	}

	// Required due to DO connection string format
	redisOptions.Username = ""

	redisClient := redis.NewClient(redisOptions)
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		logrus.Panicf("failed to connect to redis using supplied options %v %+v", err, redisOptions)
	}
	return redisClient
}
