package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

func (env *Controllers) GetGameState(c *gin.Context) {
	var request schema.GetGameStateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	stateChan, err := env.GameRepository.SubscribeGameState(request.GameId, request.PlayerId, request.PlayerKey)
	defer close(stateChan)
	if err != nil {
		logrus.Error("failed to subscribe to game state", err)
		return
	}

	c.Stream(func(w io.Writer) bool {
		if state, ok := <-stateChan; ok {
			c.SSEvent("game_state_changed", state)
			return true
		}
		return false
	})
}
