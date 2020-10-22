package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// StateUpdatesRequest represents a request to retrieve a game's state
type StateUpdatesRequest struct {
	PlayerID  models.PlayerID `form:"playerId" binding:"required"`
	GameID    models.GameID   `form:"gameId" binding:"required"`
	PlayerKey string          `form:"playerKey" binding:"required"`
}
