package models

type GameId string
type PlayerId string

//go:generate pie PlayerIds.Shuffle
type PlayerIds []PlayerId