package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// GetPlayerGameKeyRequest represents a request to retrieve a player's game key
type GetPlayerGameKeyRequest struct {
	PlayerID models.PlayerID `form:"playerId" binding:"required"`
	GameID   models.GameID   `form:"gameId" binding:"required"`
}
