package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
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
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			ctx.SSEvent("message", state)
			return true
		}
		return false
	})
}
