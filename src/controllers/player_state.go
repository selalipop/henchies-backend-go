package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (c *Controllers) GetPlayerGameKey(ctx *gin.Context) {
	var request schema.GetPlayerGameKeyRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		WriteInvalidRequestResponse(ctx, err)
		return
	}
	id, err := c.Repository.GetPlayerGameKey(ctx, request.GameId, request.PlayerId, ctx.ClientIP())
	if err != nil {
		WriteInternalErrorResponse(ctx, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (c *Controllers) GetPlayerState(ctx *gin.Context) {
	var request schema.GetPlayerStateRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerKey := models.PlayerGameKey{Key: request.PlayerKey, OwnerIp: ctx.ClientIP()}
	stateChan, err := c.Repository.SubscribePlayerState(ctx, request.GameId, request.PlayerId, playerKey)

	if err != nil {
		if err == repository.InvalidPlayerKeyErr {
			WriteAuthenticationErrorResponse(ctx, err)
		} else {
			WriteInternalErrorResponse(ctx, err)
		}
		return
	}

	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			ctx.SSEvent("player_state_changed", state)
			return true
		}
		return false
	})
}
