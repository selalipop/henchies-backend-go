package models

// PlayerColor represents the in-game color of a player
//go:generate jsonenums -type=PlayerColor
type PlayerColor int

// PlayerColors represents a list of PlayerColors
//go:generate pie PlayerColors.DropTop.Shuffle.FindFirstUsing
type PlayerColors []PlayerColor

// Color Values
const (
	Teal PlayerColor = iota
	Blue
	Amber
	Red
	Lime
	Purple
	Pink
)

// GetSelectableColors returns a list of valid in-game colors for a player
func GetSelectableColors() PlayerColors {
	return []PlayerColor{
		Teal,
		Blue,
		Amber,
		Red,
		Lime,
		Purple,
		Pink,
	}
}
