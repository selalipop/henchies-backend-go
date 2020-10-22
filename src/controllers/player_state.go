package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
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
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Writer.WriteHeader(http.StatusOK)

	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			ctx.SSEvent("message", state)
			return true
		}
		return false
	})
}
