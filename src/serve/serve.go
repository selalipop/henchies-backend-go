package main

import (
	. "github.com/SelaliAdobor/henchies-backend-go/src/controllers"
	. "github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
)

type Arguments struct {
	RedisConnectUrl     string `required:"true"`
}

func main() {
	r := gin.New()

	args := GetArguments()
	redisOptions := redis.Options{Addr: args.RedisConnectUrl}
	redisClient :=  redis.NewClient(&redisOptions)
	repositoryEnv := RepositoryEnv{
		RedisClient: redisClient,
		Context:     nil,
	}
	playerRepository := PlayerRepository{repositoryEnv}
	gameRepository := GameRepository{&playerRepository, repositoryEnv}

	controllers := Controllers{
		PlayerRepository: playerRepository,
		GameRepository:   gameRepository,
	}
	r.GET("player/state", controllers.GetPlayerState)
	r.GET("player/key", controllers.GetPlayerGameKey)

	r.GET("photon-webhooks/room-created", controllers.RoomCreatedWebhook)
	r.GET("photon-webhooks/player-joined", controllers.PlayerJoinedWebhook)
	err := r.Run()
	if(err != nil){
		logrus.Fatal(err)
	}
}

func GetArguments() Arguments {
	var args Arguments
	err := envconfig.Process("henchies", &args)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return args
}
