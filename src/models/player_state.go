package models

type PlayerGameKey struct {
	Key string `json:"key"`
	OwnerIp      string  `json:"ip"`
}

type PlayerState struct {
	GlobalUserId PlayerId      `json:"playerId"`
	GameKey      PlayerGameKey `json:"gameKey"`
	CurrentGame  GameId        `json:"currentGame"`
	IsImposter   bool          `json:"isImposter"`
	Color        PlayerColor   `json:"color"`
}
