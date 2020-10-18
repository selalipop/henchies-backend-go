package models

//go:generate jsonenums -type=PlayerColor

type PlayerColor int
//go:generate pie PlayerColors.DropTop
type PlayerColors []PlayerColor

const (
	Teal PlayerColor = iota
	Blue
	Amber
	Red
	Lime
	Purple
	Pink
)

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
