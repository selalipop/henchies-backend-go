package models



type PlayerState struct {
	GameKey      PlayerGameKey `json:"gameKey"`
	CurrentGame  GameId        `json:"currentGame"`
	IsImposter   bool          `json:"isImposter"`
	Color        PlayerColor   `json:"color"`
}
