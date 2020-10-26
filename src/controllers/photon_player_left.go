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

	logrus.Debugf("processing player s from Photon: %+v", request)

	err := c.Repository.ClearPlayerState(ctx, request.GameID, request.PlayerID)

	if err != nil {
		writeInternalErrorResponse(ctx, err)
		return
	}
	err = c.Repository.UpdateGameState(ctx, request.GameID, func(gameState models.GameState) models.GameState {
		p := gameState.Players.GetPlayerByID(request.PlayerID)
		if p == nil {
			return gameState
		}

		gameState.Players = gameState.Players.Filter(func(p models.GameStatePlayer) bool {
			return p.PlayerID != request.PlayerID
		})
		return gameState
	})

	writeSuccessIfNoErrors(ctx, err)
}
