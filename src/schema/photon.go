package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

// PhotonArgs represents base set of arguments sent by Photon to webhooks
type PhotonArgs struct {
	AppID      string        `json:"AppId"`
	AppVersion string        `json:"AppVersion"`
	Region     string        `json:"Region"`
	GameID     models.GameID `json:"GameID"`
	Type       string        `json:"Type"`
}

// PhotonExtendedArgs represents extended set of arguments sent by Photon to all webhooks besides RoomClosed
type PhotonExtendedArgs struct {
	ActorNr  string          `json:"ActorNr"`
	PlayerID models.PlayerID `json:"UserId"`
	NickName string          `json:"NickName"`
	PhotonArgs
}

// CustomRoomProperties represents custom properties sent during room creation
type CustomRoomProperties struct {
	ImposterCount int `json:"ImposterCount"`
}

// CreateOptions represents options sent during room creation
type CreateOptions struct {
	MaxPlayers         int                  `json:"MaxPlayers"`
	LobbyID            string               `json:"LobbyId"`
	LobbyType          int                  `json:"LobbyType"`
	CustomProperties   CustomRoomProperties `json:"CustomProperties"`
	EmptyRoomTTL       int                  `json:"EmptyRoomTTL"`
	PlayerTTL          int                  `json:"PlayerTTL"`
	CheckUserOnJoin    bool                 `json:"CheckUserOnJoin"`
	DeleteCacheOnLeave bool                 `json:"DeleteCacheOnLeave"`
	SuppressRoomEvents bool                 `json:"SuppressRoomEvents"`
}

// RoomCreatedRequest represents a Photon webhook call for a room being created
type RoomCreatedRequest struct {
	CreateOptions CreateOptions `json:"CreateOptions"`
	PhotonExtendedArgs
}

// PlayerJoinedRequest represents a Photon webhook call for a player joining a room
type PlayerJoinedRequest struct {
	PhotonExtendedArgs
}

// PlayerLeftRequest represents a Photon webhook call for a player leaving a room
type PlayerLeftRequest struct {
	PhotonExtendedArgs
}

// RoomClosedRequest represents a Photon webhook call for a room being closed
type RoomClosedRequest struct {
	PhotonArgs
}
