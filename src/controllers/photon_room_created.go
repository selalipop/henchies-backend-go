package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func (env *Controllers) RoomCreatedWebhook(c *gin.Context) {
	var request schema.RoomCreatedRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	var imposterCount = request.CreateOptions.CustomProperties.ImposterCount
	if imposterCount == 0 {
		//TODO: Get real imposter count
		imposterCount = int(math.Ceil( 0.2 * float64(request.CreateOptions.MaxPlayers)))
	}

	err := env.GameRepository.InitGameState(c, request.GameId, request.CreateOptions.MaxPlayers, imposterCount)
	if err != nil {
		WriteInternalErrorResponse(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
