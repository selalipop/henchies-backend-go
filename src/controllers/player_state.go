package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (env *Controllers) GetPlayerGameKey(c *gin.Context) {
	var request schema.GetPlayerGameKeyRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	id, err := env.PlayerRepository.GetPlayerGameKey(c, request.GameId, request.PlayerId, c.ClientIP())
	if err != nil {
		WriteInternalErrorResponse(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (env *Controllers) GetPlayerState(c *gin.Context) {
	var request schema.GetPlayerStateRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerKey := models.PlayerGameKey{Key: request.PlayerKey, OwnerIp: c.ClientIP()}
	stateChan, err := env.PlayerRepository.SubscribePlayerState(c, request.GameId, request.PlayerId, playerKey)

	if err != nil {
		if err == repository.InvalidPlayerKeyErr {
			WriteAuthenticationErrorResponse(c, err)
		} else {
			WriteInternalErrorResponse(c, err)
		}
		return
	}

	c.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			c.SSEvent("player_state_changed", state)
			return true
		}
		return false
	})
}
