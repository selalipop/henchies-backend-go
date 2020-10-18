package schema

import . "github.com/SelaliAdobor/henchies-backend-go/src/models"


type GetGameStateRequest struct {
	PlayerId  PlayerId      `json:"playerId" binding:"required"`
	GameId    GameId        `json:"gameId" binding:"required"`
	PlayerKey string `json:"playerKey" binding:"required"`
}