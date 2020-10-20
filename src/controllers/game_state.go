package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

func (c *Controllers) GetGameState(ctx *gin.Context) {
	var request schema.GetGameStateRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		logrus.Error("failed to parse game state request ", err)
		WriteInvalidRequestResponse(ctx, err)
		return
	}

	stateChan, err := c.Repository.SubscribeGameState(ctx, request.GameId, request.PlayerId, models.PlayerGameKey{Key: request.PlayerKey, OwnerIp: ctx.ClientIP()})

	if err != nil {
		logrus.Error("failed to subscribe to game state", err)
		return
	}

	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			ctx.SSEvent("game_state_changed", state)
			return true
		}
		close(stateChan)
		return false
	})
}
