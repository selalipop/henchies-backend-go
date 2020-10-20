package repository

// UnauthorizedPlayer is an error returned when a player's Game Key is not valid
const UnauthorizedPlayer UnauthorizedPlayerError = "player key not authorized for action"

// UnauthorizedPlayerError the type of an error returned when a player's Game Key is not valid
type UnauthorizedPlayerError string

func (e UnauthorizedPlayerError) Error() string {
	return string(e)
}
