package controllers

import (
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/SelaliAdobor/henchies-backend-go/src/websocketutils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetPlayerGameKey returns a game specific key for the player
// All subsequent calls from the player must be from the original IP,
// so this value should be stored per-game on the client
func (c *Controllers) GetPlayerGameKey(ctx *gin.Context) {
	var request schema.GetPlayerGameKeyRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		writeInvalidRequestResponse(ctx, err)
		return
	}
	id, err := c.Repository.GetPlayerGameKey(ctx, request.GameID, request.PlayerID, ctx.ClientIP())
	if err != nil {
		logrus.Errorf("failed to get player game key: %v", err)
		writeInternalErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

var playerStateUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// GetPlayerState returns a SSE stream of Player State Changes
func (c *Controllers) GetPlayerState(ctx *gin.Context) {
	var request schema.GetPlayerStateRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerKey := models.PlayerGameKey{Key: request.PlayerKey, OwnerIP: ctx.ClientIP()}
	stateChan, err := c.Repository.SubscribePlayerState(ctx, request.GameID, request.PlayerID, playerKey)

	if err != nil {
		if err == repository.UnauthorizedPlayer {
			writeAuthenticationErrorResponse(ctx, err)
		} else {
			writeInternalErrorResponse(ctx, err)
		}
		return
	}
	conn, err := playerStateUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
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
