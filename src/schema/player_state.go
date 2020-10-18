package schema

import . "github.com/SelaliAdobor/henchies-backend-go/src/models"

type GetPlayerGameKeyRequest struct {
	PlayerId PlayerId `json:"gameId" binding:"required"`
	GameId   GameId   `json:"playerId" binding:"required"`
}

type GetPlayerGameKeyResponse struct {
	PlayerKey PlayerGameKey `json:"playerKey" binding:"required"`
}

type GetPlayerStateRequest struct {
	PlayerId  PlayerId      `json:"gameId" binding:"required"`
	GameId    GameId        `json:"playerId" binding:"required"`
	PlayerKey PlayerGameKey `json:"playerKey" binding:"required"`
}

type GetPlayerStateResponse struct {
	PlayerKey PlayerGameKey `json:"playerKey" binding:"required"`
}