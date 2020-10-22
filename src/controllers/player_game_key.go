package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
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
	key, err := c.Repository.GetPlayerGameKey(ctx, request.GameID, request.PlayerID, ctx.ClientIP())
	if err != nil {
		logrus.Errorf("failed to get player game key: %v", err)
		writeInternalErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, key)
}
