package controllers

import (
	"github.com/SelaliAdobor/henchies-backend-go/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (env *Controllers) GetGameState(c *gin.Context) {
	GetGameStateWSHandler(env, c.Writer, c.Request)
}

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetGameStateWSHandler(env *Controllers, w http.ResponseWriter, r *http.Request) {
	conn, err := socketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		logrus.Error("Failed to set websocket upgrade", err)
		return
	}

	var request schema.GetGameStateRequest

	if err := conn.ReadJSON(&request); err != nil {
		logrus.Error("Failed to set game state request", err)
		_ = conn.WriteJSON(map[string]interface{}{"error" : err.Error()})
		return
	}

	stateChan, err := env.GameRepository.SubscribeGameState(request.GameId, request.PlayerId, request.PlayerKey)

	if err != nil {
		logrus.Error("Failed to subscribe to game state", err)
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
