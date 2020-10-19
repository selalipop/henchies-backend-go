package repository

import (
	"encoding/json"
	"errors"
	_ "errors"
	"fmt"
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

const PlayerIdTTl = 7200 * time.Second
const PlayerStateIdTTl = 7200 * time.Second

type PlayerRepository struct {
	RepositoryEnv
}

// Returns the game key for a given player for a specific game
// The IP address must match the first IP address used to access this key
// If the IP used to access to endpoint changes, an error is returned
// This prevents other players from using a player's game key
func (r PlayerRepository) GetPlayerGameKey(gameId GameId, playerId PlayerId, ipAddress string) (playerKey PlayerGameKey, err error) {
	var keyName = GetPlayerGameKey( gameId, playerId)
	newPlayerId := uuid.New()
	newKey := PlayerGameKey{Key: newPlayerId.String(), OwnerIp: ipAddress}

	serializedKey, err := json.Marshal(newKey)
	if err != nil {
		logrus.Errorf("Failed to marshal new player key", err.Error())
		return playerKey, err
	}

	_, err = r.RedisClient.SetNX(r.Context, keyName, serializedKey, PlayerIdTTl).Result()
	if err != nil {
		logrus.Errorf("Failed to setnx player game key from Redis", err.Error())
		return playerKey, err
	}

	result, err := r.RedisClient.Get(r.Context, keyName).Result()
	if err != nil {
		logrus.Errorf("Failed to get player game key from Redis", err.Error())
		return playerKey, err
	}
	var currentKey PlayerGameKey
	err = json.Unmarshal([]byte(result), &currentKey)
	if err != nil {
		return playerKey, err
	}

	if currentKey.OwnerIp != ipAddress {
		return playerKey, errors.New("ip address mismatch")
	}
	return currentKey, nil
}

func GetPlayerStateKey(gameId GameId, playerId PlayerId) string {
	return fmt.Sprintf("playerState:%s:%s", gameId, playerId)
}

func GetPlayerGameKey(gameId GameId, playerId PlayerId) string {
	return fmt.Sprintf("playerGameKey:%s:%s", gameId, playerId)
}

// Compare the given key to the one stored for the player
// If there is no error and the 'valid' is false, this was an attempt to access player state with the wrong key
func (env RepositoryEnv) CheckPlayerKey(gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (valid bool, err error) {
	storedKey, err := env.RedisClient.Get(env.Context, GetPlayerGameKey(gameId, playerId)).Result()
	if err != nil {
		return false, err
	}
	var key PlayerGameKey
	err = json.Unmarshal([]byte(storedKey), &key)
	if err != nil {
		return false, err
	}

	return key.Key == playerKey.Key, nil
}

//Retrieve current player state, checks for a valid player key
func (r PlayerRepository) GetPlayerStateChecked(gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (state PlayerState, err error) {
	return GetPlayerState(gameId, playerId, playerKey, true, r.RepositoryEnv)
}

//Retrieve current player state without checking ID, for internal use only, do not allow players to access unchecked player state
func (r PlayerRepository) GetPlayerStateUnchecked(gameId GameId, playerId PlayerId) (state PlayerState, err error) {
	return GetPlayerState(gameId, playerId, PlayerGameKey{}, false, r.RepositoryEnv)
}

func GetPlayerState(gameId GameId, playerId PlayerId, playerKey PlayerGameKey, shouldCheck bool, env RepositoryEnv) (state PlayerState, err error) {
	state = PlayerState{}

	if shouldCheck {
		isValidKey, err := env.CheckPlayerKey(gameId, playerId, playerKey)
		if err != nil {
			return state, err
		}
		if !isValidKey {
			return state, InvalidPlayerKeyErr
		}
	}

	var keyName = GetPlayerStateKey(gameId, playerId)

	result, err := env.RedisClient.Get(env.Context, keyName).Result()
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
func (r PlayerRepository) UpdatePlayerStateChecked(gameId GameId, playerId PlayerId, playerKey PlayerGameKey, update func(state *PlayerState)) error {
	return UpdatePlayerState(gameId, playerId, playerKey, true, r.RepositoryEnv, update)
}

//Update current player state without checking ID
//For internal use only, do not allow players to access unchecked player state
func (r PlayerRepository) UpdatePlayerStateUnchecked(gameId GameId, playerId PlayerId, update func(state *PlayerState)) error {
	return UpdatePlayerState(gameId, playerId, PlayerGameKey{}, false, r.RepositoryEnv, update)
}

func UpdatePlayerState(gameId GameId, playerId PlayerId, playerKey PlayerGameKey, shouldCheck bool, env RepositoryEnv, update func(state *PlayerState)) error {
	if shouldCheck {
		isValidKey, err := env.CheckPlayerKey(gameId, playerId, playerKey)
		if err != nil {
			return err
		}
		if !isValidKey {
			return InvalidPlayerKeyErr
		}
	}
	var keyName = GetPlayerStateKey(gameId, playerId)

	//Todo: Does this need to lock? Client is single-threaded currently, and player data is client-scoped
	state, err := GetPlayerState(gameId, playerId, playerKey, shouldCheck, env)
	if err != nil {
		return err
	}
	update(&state)

	serializedState, err := json.Marshal(state)
	if err != nil {
		logrus.Errorf("Failed to serialize player state", err.Error())
		return err
	}

	_, err = env.RedisClient.Set(env.Context, keyName, serializedState, PlayerStateIdTTl).Result()
	if err != nil {
		logrus.Errorf("Failed to save player state to Redis", err.Error())
		return err
	}
	return nil
}
