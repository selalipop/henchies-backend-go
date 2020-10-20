package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// GetPlayerGameKeyRequest represents a request to retrieve a player's game key
type GetPlayerGameKeyRequest struct {
	PlayerID models.PlayerID `json:"gameId" binding:"required"`
	GameID   models.GameID   `json:"playerId" binding:"required"`
}

// GetPlayerGameKeyResponse represents a response to GetPlayerGameKeyRequest
type GetPlayerGameKeyResponse struct {
	PlayerKey string `json:"playerKey" binding:"required"`
}

// GetPlayerStateRequest represents a request to retrieve a player's state
type GetPlayerStateRequest struct {
	PlayerID  models.PlayerID `json:"gameId" binding:"required"`
	GameID    models.GameID   `json:"playerId" binding:"required"`
	PlayerKey string          `json:"playerKey" binding:"required"`
}