package models

type PlayerId string

//go:generate pie PlayerIds.Shuffle.Contains.Filter
type PlayerIds []PlayerId
