package controllers

import (
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func (env *Controllers) RoomCreatedWebhook(c *gin.Context) {
	var request schema.RoomCreatedRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	var imposterCount = request.CustomProperties.ImposterCount
	if imposterCount == 0 {
		//TODO: Get real imposter count
		imposterCount = int(math.Ceil( 0.2 * float64(request.MaxPlayers)))
	}

	err := env.GameRepository.InitGameState(request.GameId, request.MaxPlayers, imposterCount)
	if err != nil {
		WriteInternalErrorResponse(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (env *Controllers) PlayerJoinedWebhook(c *gin.Context) {
	var request schema.PlayerJoinedRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}

	err := env.GameRepository.UpdateGameState(request.GameId, func(gameState GameState) GameState{
		gameState.Players = append(gameState.Players, request.UserId)

		if len(gameState.Players) == gameState.MaxPlayers {
			startGame(request.GameId, &gameState, env)
		}

		return gameState
	})

	if err != nil {
		WriteInternalErrorResponse(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func startGame(gameId GameId, gameState *GameState, env *Controllers) {
	gameState.Phase = Starting
	randSource := rand.NewSource(time.Now().UnixNano())

	//Shuffles player list and takes top X players as imposters
	//Also assigns a color TODO: Accept preferred colors from in-game preferences
	gameState.Players = gameState.Players.Shuffle(randSource)

	remainingColors := GetSelectableColors()

	for index, playerId := range gameState.Players {
		err := env.PlayerRepository.UpdatePlayerStateUnchecked(gameId, playerId, func(state *PlayerState) {
			state.IsImposter =  index < gameState.ImposterCount
			state.Color = remainingColors[0]
		})
		logrus.Error("failed to update player state", err)
		remainingColors = remainingColors.DropTop(1)
	}

	//Shuffle again so that player list doesn't give away imposters
	gameState.Players = gameState.Players.Shuffle(randSource)
}