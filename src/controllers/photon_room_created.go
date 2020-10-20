package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
)
// RoomCreatedWebhook is called by Photon during a room being created
func (c *Controllers) RoomCreatedWebhook(ctx *gin.Context) {
	var request schema.RoomCreatedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeInvalidRequestResponse(ctx, err)
		return
	}

	logrus.Debugf("processing room created event from Photon: %v", request)

	var imposterCount = request.CreateOptions.CustomProperties.ImposterCount
	if imposterCount == 0 {
		//TODO: Get real imposter count
		imposterCount = int(math.Ceil( 0.2 * float64(request.CreateOptions.MaxPlayers)))
	}

	err := c.Repository.InitGameState(ctx, request.GameID, request.CreateOptions.MaxPlayers, imposterCount)
	if err != nil {
		writeInternalErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
