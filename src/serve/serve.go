package main

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/controllers"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

type Arguments struct {
	RedisConnectURL string `required:"true"`
}

func main() {
	args := GetArguments()

	redisOptions, err := redis.ParseURL(args.RedisConnectURL)
	if err != nil {
		logrus.Error("failed to parse Redis Connection Url ", err)
	}

	// Required due to DO connection string format
	redisOptions.Username = ""

	redisClient := redis.NewClient(redisOptions)

	c := controllers.Controllers{
		Repository: repository.Repository{
			RedisClient: redisClient,
		},
	}

	g := gin.New()

	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.TraceLevel)
	g.Use(ginlogrus.Logger(logger), gin.Recovery())

	SetupRoutes(g, c)

	err = g.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}

//goland:noinspection ALL
func SetupRoutes(g *gin.Engine, c controllers.Controllers) {
	g.GET("/", c.GetInfo)

	g.GET("/player/state", c.GetPlayerState)
	g.GET("/player/key", c.GetPlayerGameKey)

	g.GET("/game/state", c.GetGameState)

	g.POST("/photonwebhooks/roomcreated", c.RoomCreatedWebhook)
	g.POST("/photonwebhooks/roomclosed", c.RoomClosedWebhook)

	g.POST("/photonwebhooks/playerjoined", c.PlayerJoinedWebhook)
	g.POST("/photonwebhooks/playerleft", c.PlayerLeftWebhook)
}

func GetArguments() Arguments {
	var args Arguments
	const argumentsPrefix = "henchies"
	err := envconfig.Process(argumentsPrefix, &args)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return args
}
