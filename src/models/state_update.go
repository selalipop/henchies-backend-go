package models

// StateUpdate represents the current in-game state of a game
type StateUpdate struct {
	PlayerState *PlayerState `json:"playerState"`
	GameState   *GameState   `json:"gameState"`
}

// ToUpdate converts state to an update
func (s *PlayerState) ToUpdate() StateUpdate {
	return StateUpdate{
		PlayerState: s,
		GameState:   nil,
	}
}

// ToUpdate converts state to an update
func (s *GameState) ToUpdate() StateUpdate {
	return StateUpdate{
		PlayerState: nil,
		GameState:   s,
	}
}
