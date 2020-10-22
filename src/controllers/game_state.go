package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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

	resp := ctx.Writer
	h := resp.Header()
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")
	h.Set("X-Accel-Buffering", "no")
	h.Set("Access-Control-Allow-Origin", "*")

	resp.WriteHeader(http.StatusOK)
	resp.WriteHeaderNow()
	resp.Flush()
	logrus.Trace("wrote headers")

	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			logrus.Trace("writing message")
			ctx.SSEvent("message", state)
			logrus.Trace("wrote message")
			return true
		}
		return false
	})
}
