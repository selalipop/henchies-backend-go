package models

// GameState represents the current in-game state of a game
type GameState struct {
	MaxPlayers    int `json:"maxPlayers"`
	ImposterCount int `json:"imposterCount"`
	Phase         GamePhase `json:"phase"`
	Players       PlayerIDs `json:"players"`
}
