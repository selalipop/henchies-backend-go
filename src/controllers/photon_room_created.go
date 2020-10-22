package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
)

const gameLobbyId = "GameLobby"

// RoomCreatedWebhook is called by Photon during a room being created
func (c *Controllers) RoomCreatedWebhook(ctx *gin.Context) {
	var request schema.RoomCreatedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		logrus.Debugf("failed to parse room created event from Photon: %v", err)
		writeInvalidRequestResponse(ctx, err)
		return
	}

	logrus.Debugf("processing room created event from Photon: %+v", request)

	if request.CreateOptions.LobbyID != gameLobbyId {
		logrus.Debugf("ignoring room created outside of game lobby: %+v", request)
		writeSuccessIfNoErrors(ctx)
	}

	var imposterCount = request.CreateOptions.CustomProperties.ImposterCount
	if imposterCount == 0 {
		//TODO: Get real imposter count
		imposterCount = int(math.Ceil(0.2 * float64(request.CreateOptions.MaxPlayers)))
	}

	err := c.Repository.InitGameState(ctx, request.GameID, request.CreateOptions.MaxPlayers, imposterCount)
	if err != nil {
		logrus.Errorf("failed to initialize game state on room created event: %+v err: %+v", request, err)
		writeInternalErrorResponse(ctx, err)
	}

	if !request.CreateOptions.CustomProperties.ServerCreatedRoom {
		// Photon Treats the RoomCreated webhook as a proxy for PlayerJoined in some cases
		c.processPlayerJoined(ctx, request.GameID, request.PlayerID)
	}
}
