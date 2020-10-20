package controllers

import (
	"context"
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const WaitForLeavingPlayersDuration = 15 * time.Second

func (env *Controllers) PlayerJoinedWebhook(c *gin.Context) {
	var request schema.PlayerJoinedRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	err := env.PlayerRepository.UpdatePlayerStateUnchecked(c,request.GameId, request.UserId, func(state PlayerState) PlayerState {
		state.CurrentGame = request.GameId
		return state
	})

	if err != nil {
		WriteInternalErrorResponse(c, err)
		return
	}
	err = env.GameRepository.UpdateGameState(c, request.GameId, func(gameState GameState) GameState {
		if gameState.Players.Contains(request.UserId) {
			return gameState
		}

		gameState.Players = append(gameState.Players, request.UserId)

		if len(gameState.Players) == gameState.MaxPlayers {
			gameState.Phase = Starting
			go startGame(c, request.GameId, env)
		}
		return gameState
	})

	if err != nil {
		WriteInternalErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func startGame(ctx context.Context, gameId GameId, env *Controllers) {
	err := env.GameRepository.UpdateGameState(ctx, gameId, func(gameState GameState) GameState {
		if gameState.Phase != Starting {
			return gameState
		}

		if len(gameState.Players) < gameState.MaxPlayers {
			gameState.Phase = WaitingForPlayers
			return gameState
		}

		time.Sleep(WaitForLeavingPlayersDuration)

		randSource := rand.NewSource(time.Now().UnixNano())

		//Shuffles player list and takes top X players as imposters
		//Also assigns a color TODO: Accept preferred colors from in-game preferences
		gameState.Players = gameState.Players.Shuffle(randSource)

		remainingColors := GetSelectableColors()

		for index, playerId := range gameState.Players {
			err := env.PlayerRepository.UpdatePlayerStateUnchecked(ctx, gameId, playerId, func(state PlayerState) PlayerState {
				state.IsImposter = index < gameState.ImposterCount
				state.Color = remainingColors[0]
				return state
			})
			if err != nil {
				logrus.Error("failed to update player state", err)
			}
			remainingColors = remainingColors.DropTop(1)
		}

		//Shuffle again so that player list doesn't give away imposters
		gameState.Players = gameState.Players.Shuffle(randSource)
		gameState.Phase = Started
		return gameState
	})

	if err != nil {
		logrus.Error("failed to start game ", err)
	}
}
