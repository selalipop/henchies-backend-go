package controllers

import (
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controllers) PlayerLeftWebhook(ctx *gin.Context) {
	var request schema.PlayerLeftRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(ctx, err)
		return
	}

	err := c.Repository.UpdatePlayerStateUnchecked(ctx, request.GameId, request.UserId, func(state PlayerState) PlayerState {
		if state.CurrentGame == request.GameId {
			state = PlayerState{}
		}
		return state
	})

	if err != nil {
		WriteInternalErrorResponse(ctx, err)
		return
	}
	err = c.Repository.UpdateGameState(ctx, request.GameId, func(gameState GameState) GameState {
		if !gameState.Players.Contains(request.UserId) {
			return gameState
		}

		gameState.Players = gameState.Players.Filter(func(playerId PlayerId) bool {
			return playerId != request.UserId
		})
		return gameState
	})

	if err != nil {
		WriteInternalErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
