package models

// PlayerGameKey is used to authenticate player actions and access to various state
type PlayerGameKey struct {
	Key     string `json:"key"`
	OwnerIP string `json:"ip"`
}