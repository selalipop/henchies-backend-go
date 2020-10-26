package models

import (
	"math/rand"
	"time"
)

// GameStatePlayer represents game-scoped data specific to each player
type GameStatePlayer struct {
	PlayerID    PlayerID    `json:"playerId"`
	PlayerColor PlayerColor `json:"playerColor"`
}

// GameStatePlayers represents a list of GameStatePlayers
//go:generate pie GameStatePlayers.Shuffle.Contains.Filter.FindFirstUsing
type GameStatePlayers []GameStatePlayer

func (p GameStatePlayers) GetUnusedColor() PlayerColor {
	randSource := rand.NewSource(time.Now().UnixNano())
	colors := GetSelectableColors().Shuffle(randSource)
	unusedColor := colors.FindFirstUsing(func(color PlayerColor) bool {
		return p.FindFirstUsing(func(p GameStatePlayer) bool {
			return p.PlayerColor == color
		}) == -1
	})
	return colors[unusedColor]
}

func (p GameStatePlayers) GetPlayerByID(id PlayerID) *GameStatePlayer {
	for _, value := range p {
		if value.PlayerID == id {
			return &value
		}
	}

	return nil
}
