package models

// GameState represents the current state of a game
type GameState struct {
	MaxPlayers    int              `json:"maxPlayers"`
	ImposterCount int              `json:"imposterCount"`
	Phase         GamePhase        `json:"phase"`
	Players       GameStatePlayers `json:"players"`
}
