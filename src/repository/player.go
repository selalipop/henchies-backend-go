package repository

import (
	. "context"
	"encoding/json"
	"errors"
	_ "errors"
	"fmt"
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/redisUtil"
	. "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

const PlayerIdTTl = 7200 * time.Second
const PlayerStateIdTTl = 7200 * time.Second

type PlayerRepository struct {
	Repository
}

// Returns the game key for a given player for a specific game
// The IP address must match the first IP address used to access this key
// If the IP used to access to endpoint changes, an error is returned
// This prevents other players from using a player's game key
func (r PlayerRepository) GetPlayerGameKey(ctx Context, gameId GameId, playerId PlayerId, ipAddress string) (playerKey PlayerGameKey, err error) {
	var keyName = GetPlayerGameKey(gameId, playerId)
	newPlayerId := uuid.New()
	newKey := PlayerGameKey{Key: newPlayerId.String(), OwnerIp: ipAddress}

	serializedKey, err := json.Marshal(newKey)

	if err != nil {
		logrus.Errorf("Failed to marshal new player key", err.Error())
		return playerKey, err
	}

	_, err = r.RedisClient.SetNX(ctx, keyName, serializedKey, PlayerIdTTl).Result()
	if err != nil {
		logrus.Errorf("Failed to setnx player game key from Redis", err.Error())
		return playerKey, err
	}

	result, err := r.RedisClient.Get(ctx, keyName).Result()
	if err != nil {
		logrus.Errorf("Failed to get player game key from Redis", err.Error())
		return playerKey, err
	}
	var currentKey PlayerGameKey
	err = json.Unmarshal([]byte(result), &currentKey)
	if err != nil {
		return playerKey, err
	}

	if currentKey == newKey {
		err := r.UpdatePlayerStateUnchecked(ctx, gameId, playerId, func(state PlayerState) PlayerState {
			state.GameKey = currentKey
			return state
		})
		if err != nil{
			return playerKey, err
		}
	}

	if currentKey.OwnerIp != ipAddress {
		return playerKey, errors.New("ip address mismatch")
	}
	return currentKey, nil
}

func GetPlayerStateKey(gameId GameId, playerId PlayerId) string {
	return fmt.Sprintf("playerState:%s:%s", gameId, playerId)
}
func GetPlayerStatePublishKey(gameId GameId, playerId PlayerId) string {
	return fmt.Sprintf("playerPublishState:%s:%s", gameId, playerId)
}
func GetPlayerGameKey(gameId GameId, playerId PlayerId) string {
	return fmt.Sprintf("playerGameKey:%s:%s", gameId, playerId)
}

// Compare the given key to the one stored for the player
// If there is no error and the 'valid' is false, this was an attempt to access player state with the wrong key
func (env Repository) CheckPlayerKey(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (valid bool, err error) {

	storedKeyJson, err := env.RedisClient.Get(ctx, GetPlayerGameKey(gameId, playerId)).Result()
	if err != nil {
		return false, err
	}
	var storedKey PlayerGameKey
	err = json.Unmarshal([]byte(storedKeyJson), &storedKey)
	if err != nil {
		return false, err
	}

	return storedKey.Key == playerKey.Key, nil
}

//Retrieve current player state, checks for a valid player key
func (r PlayerRepository) GetPlayerStateChecked(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (state PlayerState, err error) {

	return GetPlayerState(ctx, gameId, playerId, playerKey, true, r.Repository)
}

//Retrieve current player state without checking ID, for internal use only, do not allow players to access unchecked player state
func (r PlayerRepository) GetPlayerStateUnchecked(ctx Context,
	gameId GameId, playerId PlayerId) (state PlayerState, err error) {

	return GetPlayerState(ctx, gameId, playerId, PlayerGameKey{}, false, r.Repository)
}

func GetPlayerState(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey, shouldCheck bool, env Repository) (state PlayerState, err error) {

	state = PlayerState{}

	if shouldCheck {
		isValidKey, err := env.CheckPlayerKey(ctx, gameId, playerId, playerKey)
		if err != nil {
			return state, err
		}
		if !isValidKey {
			return state, InvalidPlayerKeyErr
		}
	}

	var keyName = GetPlayerStateKey(gameId, playerId)

	result, err := env.RedisClient.Get(ctx, keyName).Result()
	if err != nil {
		logrus.Errorf("Failed to get player state from Redis", err.Error())
		return state, err
	}
	err = json.Unmarshal([]byte(result), &state)

	if err != nil {
		logrus.Errorf("Failed to deserialize player state from Redis", err.Error())
		return state, err
	}
	return state, nil
}

//Update current player state while checking ID
//Returns InvalidPlayerKeyErr if the key passed did not match the given player ID
func (r PlayerRepository) UpdatePlayerStateChecked(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey, update func(state PlayerState) PlayerState) error {
	return UpdatePlayerState(ctx, gameId, playerId, playerKey, true, r.Repository, update)
}

//Update current player state without checking ID
//For internal use only, do not allow players to access unchecked player state
func (r PlayerRepository) UpdatePlayerStateUnchecked(ctx Context,
	gameId GameId, playerId PlayerId, update func(state PlayerState) PlayerState) error {

	return UpdatePlayerState(ctx, gameId, playerId, PlayerGameKey{}, false, r.Repository, update)
}

func UpdatePlayerState(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey, shouldCheck bool, env Repository, update func(state PlayerState) PlayerState) error {

	if shouldCheck {
		isValidKey, err := env.CheckPlayerKey(ctx, gameId, playerId, playerKey)
		if err != nil {
			return err
		}
		if !isValidKey {
			return InvalidPlayerKeyErr
		}
	}

	return UpdatePlayerStateTransaction(ctx, env.RedisClient, gameId, playerId, update)
}

func UpdatePlayerStateTransaction(ctx Context,
	client *Client, gameId GameId, playerId PlayerId, update func(state PlayerState) PlayerState) error {

	stateKey := GetPlayerStateKey(gameId, playerId)
	publishKey := GetPlayerStatePublishKey(gameId, playerId)

	return redisUtil.UpdateKeyTransaction(ctx, client, stateKey, publishKey, GameStateTTL, 0,
		func(data []byte) (interface{}, error) {
			var playerState PlayerState
			err := json.Unmarshal(data, &playerState)
			return playerState, err
		},
		func() interface{} {
			return PlayerState{}
		}, func(value interface{}) interface{} {
			return update(value.(PlayerState))
		})
}

func (r *PlayerRepository) SubscribePlayerState(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (channel chan PlayerState, err error) {

	valid, err := r.CheckPlayerKey(ctx, gameId, playerId, playerKey)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidPlayerKeyErr
	}

	var playerState PlayerState
	subscription, err := redisUtil.SubscribeJson(ctx, r.RedisClient,GetPlayerStateKey(gameId,playerId), GetPlayerStatePublishKey(gameId, playerId), &playerState)

	channel = make(chan PlayerState)
	go func() {
		defer close(channel)
		for {
			latest, ok := <-subscription
			if !ok {
				return
			}
			channel <- *latest.(*PlayerState)
		}
	}()
	return channel, nil
}
