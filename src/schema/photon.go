package schema

import "github.com/SelaliAdobor/henchies-backend-go/src/models"

type PhotonArgs struct {
	AppId string `json:"AppId"`
	AppVersion      string  `json:"AppVersion"`
	Region  string         `json:"Region"`
	GameId   models.GameId           `json:"GameId"`
	Type        string    `json:"Type"`
}

type PhotonExtendedArgs struct {
	ActorNr  string          `json:"ActorNr"`
	UserId   models.PlayerId `json:"UserId"`
	NickName        string   `json:"NickName"`
	PhotonArgs
}

type  CustomRoomProperties struct{
	ImposterCount  int          `json:"ImposterCount"`
}
type CreateOptions struct {
	MaxPlayers         int                  `json:"MaxPlayers"`
	LobbyId            string               `json:"LobbyId"`
	LobbyType          int                  `json:"LobbyType"`
	CustomProperties   CustomRoomProperties `json:"CustomProperties"`
	EmptyRoomTTL       int                  `json:"EmptyRoomTTL"`
	PlayerTTL          int                  `json:"PlayerTTL"`
	CheckUserOnJoin    bool                 `json:"CheckUserOnJoin"`
	DeleteCacheOnLeave bool                 `json:"DeleteCacheOnLeave"`
	SuppressRoomEvents bool                 `json:"SuppressRoomEvents"`

}

//Requests

type RoomCreatedRequest struct {
	CreateOptions         CreateOptions                  `json:"CreateOptions"`
	PhotonExtendedArgs
}

type PlayerJoinedRequest struct {
	PhotonExtendedArgs
}

type PlayerLeftRequest struct {
	PhotonExtendedArgs
}