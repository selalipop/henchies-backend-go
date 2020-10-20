package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func (c *Controllers) RoomCreatedWebhook(ctx *gin.Context) {
	var request schema.RoomCreatedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(ctx, err)
		return
	}
	var imposterCount = request.CreateOptions.CustomProperties.ImposterCount
	if imposterCount == 0 {
		//TODO: Get real imposter count
		imposterCount = int(math.Ceil( 0.2 * float64(request.CreateOptions.MaxPlayers)))
	}

	err := c.Repository.InitGameState(ctx, request.GameId, request.CreateOptions.MaxPlayers, imposterCount)
	if err != nil {
		WriteInternalErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
