package models
//go:generate jsonenums -type=GamePhase

type GamePhase int

const (
	WaitingForPlayers GamePhase = iota
	Starting
	Started
)
