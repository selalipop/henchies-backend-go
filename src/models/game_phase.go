package models

// GamePhase the current phase of gameplay
//go:generate jsonenums -type=GamePhase
type GamePhase int

const (
	// WaitingForPlayers is a state where players are able to join, max player count not reached
	WaitingForPlayers GamePhase = iota
	// Starting is a state where enough players have joined to start, but waiting in case players leave
	Starting
	// Started is a state where gameplay has started player data such as imposter status has been set
	Started
)
