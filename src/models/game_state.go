package models

type GameState struct {
	MaxPlayers int
	ImposterCount int
	Phase GamePhase
	Players PlayerIds
}
