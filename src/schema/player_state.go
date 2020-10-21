package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// GetPlayerGameKeyRequest represents a request to retrieve a player's game key
type GetPlayerGameKeyRequest struct {
	PlayerID models.PlayerID `form:"playerId" binding:"required"`
	GameID   models.GameID   `form:"gameId" binding:"required"`
}

// GetPlayerGameKeyResponse represents a response to GetPlayerGameKeyRequest
type GetPlayerGameKeyResponse struct {
	PlayerKey string `json:"playerKey" binding:"required"`
}

// GetPlayerStateRequest represents a request to retrieve a player's state
type GetPlayerStateRequest struct {
	PlayerID  models.PlayerID `form:"playerId" binding:"required"`
	GameID    models.GameID   `form:"gameId" binding:"required"`
	PlayerKey string          `form:"playerKey" binding:"required"`
}
