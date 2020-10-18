package schema

import . "github.com/SelaliAdobor/henchies-backend-go/src/models"


type GetGameStateRequest struct {
	PlayerId  PlayerId      `json:"gameId" binding:"required"`
	GameId    GameId        `json:"playerId" binding:"required"`
	PlayerKey PlayerGameKey `json:"playerKey" binding:"required"`
}