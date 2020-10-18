package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/repository"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Controllers) GetPlayerGameKey(c *gin.Context) {
	var request schema.GetPlayerGameKeyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	id, err := env.PlayerRepository.GetPlayerGameKey(request.GameId, request.PlayerId, c.ClientIP())
	if err != nil {
		WriteInternalErrorResponse(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (env *Controllers) GetPlayerState(c *gin.Context) {
	var request schema.GetPlayerStateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state, err := env.PlayerRepository.GetPlayerStateChecked(request.GameId, request.PlayerId, request.PlayerKey)
	if err != nil {
		if err == repository.InvalidPlayerKeyErr {
			WriteAuthenticationErrorResponse(c, err)
		} else {
			WriteInternalErrorResponse(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"state": state})
}
