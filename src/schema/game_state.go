package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// GetGameStateRequest represents a request to retrieve a game's state
type GetGameStateRequest struct {
	PlayerID  models.PlayerID `form:"playerId" binding:"required"`
	GameID    models.GameID   `form:"gameId" binding:"required"`
	PlayerKey string          `form:"playerKey" binding:"required"`
}
