package models

// PlayerState represents player's in-game state
type PlayerState struct {
	GameKey     PlayerGameKey `json:"gameKey"`
	CurrentGame GameID        `json:"currentGame"`
	IsImposter  bool          `json:"isImposter"`
	Color       PlayerColor   `json:"color"`
}
