package main

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/controllers"
	"github.com/gin-gonic/gin"
)

//goland:noinspection ALL
func setupRoutes(g *gin.Engine, c controllers.Controllers) {
	g.GET("/", c.GetInfo)

	g.GET("/player/state", c.GetPlayerState)
	g.GET("/player/key", c.GetPlayerGameKey)

	g.GET("/game/state", c.GetGameState)

	g.POST("/photonwebhooks/roomcreated", c.RoomCreatedWebhook)
	g.POST("/photonwebhooks/roomclosed", c.RoomClosedWebhook)

	g.POST("/photonwebhooks/playerjoined", c.PlayerJoinedWebhook)
	g.POST("/photonwebhooks/playerleft", c.PlayerLeftWebhook)
}
