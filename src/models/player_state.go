package models

// PlayerState represents player-scoped data for the current game
type PlayerState struct {
	GameKey     PlayerGameKey `json:"gameKey"`
	CurrentGame GameID        `json:"currentGame"`
	IsImposter  bool          `json:"isImposter"`
}
