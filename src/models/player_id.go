package models

// PlayerID represents the unique ID for a Player. Should match Photon UserID
type PlayerID string

// PlayerIDs represents a list of PlayerID
//go:generate pie PlayerIds.Shuffle.Contains.Filter
type PlayerIDs []PlayerID
