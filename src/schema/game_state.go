package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// GetGameStateRequest represents a request to retrieve a game's state
type GetGameStateRequest struct {
	PlayerID  models.PlayerID `json:"playerId" binding:"required"`
	GameID    models.GameID   `json:"gameId" binding:"required"`
	PlayerKey string          `json:"playerKey" binding:"required"`
}