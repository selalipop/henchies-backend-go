package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RoomCreatedWebhook is called by Photon during a room being created
func (c *Controllers) RoomClosedWebhook(ctx *gin.Context) {
	var request schema.RoomClosedRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeInvalidRequestResponse(ctx, err)
		return
	}

	logrus.Debugf("processing room closed event from Photon: %+v", request)

	err := c.Repository.ClearGameState(ctx, request.GameID)

	writeSuccessIfNoErrors(ctx, err)
}
