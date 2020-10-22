package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/ginutil"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//  GetGameState returns a SSE stream of Game State Changes
func (c *Controllers) GetGameState(ctx *gin.Context) {
	var request schema.GetGameStateRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		logrus.Error("failed to parse game state request ", err)
		writeInvalidRequestResponse(ctx, err)
		return
	}

	stateChan, err := c.Repository.SubscribeGameState(ctx, request.GameID, request.PlayerID, models.PlayerGameKey{Key: request.PlayerKey, OwnerIP: ctx.ClientIP()})

	if err != nil {
		logrus.Error("failed to subscribe to game state", err)
		return
	}

	channel := make(chan interface{})
	go func() {
		for state := range stateChan {
			channel <- state
		}
	}()
	ginutil.ChannelToServerSentEvents(ctx, channel)
}
