package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

func (env *Controllers) GetGameState(c *gin.Context) {
	var request schema.GetGameStateRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logrus.Error("failed to parse game state request ", err)
		WriteInvalidRequestResponse(c, err)
		return
	}
	stateChan, err := env.GameRepository.SubscribeGameState(request.GameId, request.PlayerId, models.PlayerGameKey{Key: request.PlayerKey, OwnerIp: c.ClientIP()})

	if err != nil {
		logrus.Error("failed to subscribe to game state", err)
		return
	}

	c.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			c.SSEvent("game_state_changed", state)
			return true
		}
		close(stateChan)
		return false
	})
}
