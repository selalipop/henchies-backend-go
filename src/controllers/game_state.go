package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
)

func (env *Controllers) GetGameState(c *gin.Context) {
	var request schema.GetGameStateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	c.Stream(func(w io.Writer) bool {

		stateChan, err := env.GameRepository.SubscribeGameState(request.GameId, request.PlayerId, request.PlayerKey)
		defer close(stateChan)

		if err != nil {
			logrus.Error("Failed to set websocket upgrade", err)
			return false
		}

		for {
			eventId := 0

			state := <-stateChan

			err = sse.Encode(w, sse.Event{
				Id:    strconv.Itoa(eventId),
				Event: "game_state_changed",
				Data:  state,
			})

			if err != nil {
				logrus.Error("Failed to set websocket upgrade", err)
				return true
			}
		}
	})

}
