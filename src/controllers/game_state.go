package controllers

import (
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/SelaliAdobor/henchies-backend-go/src/websocketutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var gameStateUpdater = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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

	conn, err := gameStateUpdater.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		writeInternalErrorResponse(ctx, fmt.Errorf("failed to upgrade websocket: %w", err))
		return
	}

	for {
		state, ok := <-stateChan
		if !ok {
			break
		}
		err := websocketutils.WriteValueToWebsocket(state, conn)
		if err != nil {
			logrus.Errorf("failed to write update to player state socket: %v", err)
			break
		}
	}
}
