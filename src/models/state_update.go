package models

// PingUpdate represents an update sent to keep connections alive
var PingUpdate = StateUpdate{IsPing: true}

// StateUpdate represents the current in-game state of a game
type StateUpdate struct {
	PlayerState *PlayerState `json:"playerState"`
	GameState   *GameState   `json:"gameState"`
	// A State Update that is for keeping connections alive, do not send state inside this update
	IsPing bool `json:"isPing"`
}

// StateUpdateField represents a value held inside a StateUpdate
type StateUpdateField interface {
	ToUpdate() StateUpdate
}

func (s StateUpdate) ToUpdate() StateUpdate {
	return s
}

// ToUpdate converts state to an update
func (s PlayerState) ToUpdate() StateUpdate {
	return StateUpdate{
		PlayerState: &s,
		GameState:   nil,
	}
}

// ToUpdate converts state to an update
func (s GameState) ToUpdate() StateUpdate {
	return StateUpdate{
		PlayerState: nil,
		GameState:   &s,
	}
}
