package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PlayerLeftWebhook is called by Photon during a player leaving room
func (c *Controllers) PlayerLeftWebhook(ctx *gin.Context) {
	var request schema.PlayerLeftRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeInvalidRequestResponse(ctx, err)
		return
	}

	logrus.Debugf("processing player left event from Photon: %+v", request)

	err := c.Repository.ClearPlayerState(ctx, request.GameID, request.PlayerID)

	if err != nil {
		writeInternalErrorResponse(ctx, err)
		return
	}
	err = c.Repository.UpdateGameState(ctx, request.GameID, func(gameState models.GameState) models.GameState {
		if !gameState.Players.Contains(request.PlayerID) {
			return gameState
		}

		gameState.Players = gameState.Players.Filter(func(playerID models.PlayerID) bool {
			return playerID != request.PlayerID
		})
		return gameState
	})

	writeSuccessIfNoErrors(ctx, err)
}
