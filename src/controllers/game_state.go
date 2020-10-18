package controllers

import (
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (env *Controllers) GetGameState(c *gin.Context) {
	var request schema.GetGameStateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		WriteInvalidRequestResponse(c, err)
		return
	}
	GetGameStateWSHandler(env, request.GameId, request.PlayerId, request.PlayerKey, c.Writer, c.Request)
}

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetGameStateWSHandler(env *Controllers, gameId GameId, playerId PlayerId, playerKey PlayerGameKey, w http.ResponseWriter, r *http.Request) {
	conn, err := socketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		logrus.Error("Failed to set websocket upgrade", err)
		return
	}

	stateChan, err := env.GameRepository.SubscribeGameState(gameId, playerId, playerKey)
	if err != nil {
		logrus.Error("Failed to set websocket upgrade", err)
		return
	}
	for {
		state := <-stateChan
		err = conn.WriteJSON(state)
		if err != nil {
			break
		}
	}
}
