package controllers

import (
	"context"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const waitForLeavingPlayersDuration = 15 * time.Second

// PlayerJoinedWebhook is called by Photon during a player joining room
func (c *Controllers) PlayerJoinedWebhook(ctx *gin.Context) {
	var request schema.PlayerJoinedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeInvalidRequestResponse(ctx, err)
		return
	}

	logrus.Debugf("processing player joined event from Photon: %+v", request)

	c.processPlayerJoined(ctx, request.GameID, request.PlayerID)
}

func (c *Controllers) processPlayerJoined(ctx *gin.Context, gameID models.GameID, playerID models.PlayerID) {
	err := c.Repository.UpdatePlayerStateUnchecked(ctx, gameID, playerID, func(state models.PlayerState) models.PlayerState {
		state = models.PlayerState{
			CurrentGame: gameID,
		}
		return state
	})

	if err != nil {
		writeInternalErrorResponse(ctx, err)
		return
	}
	err = c.Repository.UpdateGameState(ctx, gameID, func(gameState models.GameState) models.GameState {
		if gameState.Players.GetPlayerByID(playerID) != nil {
			logrus.Debugf("received player joined but player was already in game:%+v player: %+v", gameID, playerID)
			return gameState
		}

		newPlayer := models.GameStatePlayer{
			PlayerID:    playerID,
			PlayerColor: gameState.Players.GetUnusedColor(),
		}

		gameState.Players = append(gameState.Players, newPlayer)

		if len(gameState.Players) == gameState.MaxPlayers {
			logrus.Debugf("starting game after player joined: game:%+v player: %+v", gameID, playerID)

			gameState.Phase = models.Starting
			go startGame(ctx, gameID, c)
		}
		return gameState
	})

	writeSuccessIfNoErrors(ctx, err)
}

func startGame(ctx context.Context, gameID models.GameID, env *Controllers) {
	updateMapper := func(gameState models.GameState) models.GameState {
		if gameState.Phase != models.Starting {
			return gameState
		}

		if len(gameState.Players) < gameState.MaxPlayers {
			gameState.Phase = models.WaitingForPlayers
			return gameState
		}

		time.Sleep(waitForLeavingPlayersDuration)

		randSource := rand.NewSource(time.Now().UnixNano())

		// Shuffles player list and takes top X players as imposters
		// Also assigns a color TODO: Accept preferred colors from in-game preferences
		gameState.Players = gameState.Players.Shuffle(randSource)

		for index, player := range gameState.Players {
			err := env.Repository.UpdatePlayerStateUnchecked(ctx, gameID, player.PlayerID, func(state models.PlayerState) models.PlayerState {
				state.IsImposter = index < gameState.ImposterCount
				return state
			})
			if err != nil {
				logrus.Errorf("failed to update player state %v", err)
			}
		}

		// Shuffle again so that player list doesn't give away imposters
		gameState.Players = gameState.Players.Shuffle(randSource)
		gameState.Phase = models.Started
		return gameState
	}

	err := env.Repository.UpdateGameState(ctx, gameID, updateMapper)

	if err != nil {
		logrus.Error("failed to start game ", err)
	}
}
