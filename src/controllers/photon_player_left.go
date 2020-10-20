package controllers

import (
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Controllers) PlayerLeftWebhook(c *gin.Context) {
	var request schema.PlayerLeftRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}

	err := env.PlayerRepository.UpdatePlayerStateUnchecked(c, request.GameId, request.UserId, func(state PlayerState) PlayerState {
		if state.CurrentGame == request.GameId {
			state = PlayerState{}
		}
		return state
	})

	if err != nil {
		WriteInternalErrorResponse(c, err)
		return
	}
	err = env.GameRepository.UpdateGameState(c, request.GameId, func(gameState GameState) GameState {
		if !gameState.Players.Contains(request.UserId) {
			return gameState
		}

		gameState.Players = gameState.Players.Filter(func(playerId PlayerId) bool {
			return playerId != request.UserId
		})
		return gameState
	})

	if err != nil {
		WriteInternalErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
