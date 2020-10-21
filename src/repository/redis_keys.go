package repository

import (
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"strings"
)

// Player Keys

func redisKeyPlayerGameKey(gameID models.GameID, playerID models.PlayerID) string {
	return strings.ToLower(fmt.Sprintf("player:%s:game:%s:key", playerID, gameID))
}

func redisKeyPlayerState(gameID models.GameID, playerID models.PlayerID) string {
	return strings.ToLower(fmt.Sprintf("player:%s:game:%s:state:current", playerID, gameID))
}

func redisKeyPlayerStatePublish(gameID models.GameID, playerID models.PlayerID) string {
	return strings.ToLower(fmt.Sprintf("player:%s:game:%s:state:pubSub", playerID, gameID))
}

// Game Keys

func redisPublishKeyGameState(gameID models.GameID) string {
	return strings.ToLower(fmt.Sprintf("game:%s:state:pubSub", gameID))
}

func redisKeyGameState(gameID models.GameID) string {
	return strings.ToLower(fmt.Sprintf("game:%s:state:current", gameID))
}
