package repository

import (
	. "context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/redisUtil"
	. "github.com/cenkalti/backoff"
	. "github.com/go-redis/redis/v8"
	"time"
)


func GetGameStatePubSubKey(gameId GameId) string {
	return fmt.Sprintf("playerGamePubSubKey:%s", gameId)
}
func GetGameStateKey(gameId GameId) string {
	return fmt.Sprintf("playerGameKey:%s", gameId)
}

var GameStateTTL = 5 * time.Hour
var MaxElapsedTimeGameStateUpdate = 5 * time.Minute

func (r *Repository) InitGameState(ctx Context, gameId GameId, startingPlayerCount int, imposterCount int) error {
	exists, err :=
		r.RedisClient.Exists(ctx, GetGameStateKey(gameId)).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	return UpdateGameStateTransaction(ctx, r.RedisClient, gameId, func(gameState GameState) GameState {
		return GameState{
			ImposterCount: imposterCount,
			MaxPlayers:    startingPlayerCount,
			Phase:         WaitingForPlayers,
			Players:       []PlayerId{},
		}
	})
}
func (r *Repository) AddPlayerToGame(ctx Context,
	gameId GameId, playerId PlayerId) error {

	exists, err :=
		r.RedisClient.Exists(ctx, GetGameStateKey(gameId)).Result()

	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	err = UpdateGameStateTransaction(ctx, r.RedisClient, gameId, func(gameState GameState) GameState {
		gameState.Players = append(gameState.Players, playerId)
		return gameState
	})
	if err != nil {
		return err
	}

	return r.UpdatePlayerStateUnchecked(ctx, gameId, playerId, func(state PlayerState) PlayerState {
		state.CurrentGame = gameId
		return state
	})
}
func (r *Repository) UpdateGameState(ctx Context,
	gameId GameId, update func(gameState GameState) GameState) error {

	operation := func() error {
		return UpdateGameStateTransaction(ctx, r.RedisClient, gameId, update)
	}

	backoff := NewExponentialBackOff()
	backoff.MaxElapsedTime = MaxElapsedTimeGameStateUpdate
	return Retry(operation, backoff)
}

func UpdateGameStateTransaction(ctx Context,
	client *Client, gameId GameId, update func(gameState GameState) GameState) error {

	stateKey := GetGameStateKey(gameId)
	publishKey := GetGameStatePubSubKey(gameId)

	return redisUtil.UpdateKeyTransaction(ctx, client, stateKey, publishKey, GameStateTTL, 0,
		func(data []byte) (interface{}, error) {
			var gameState GameState
			err := json.Unmarshal(data, &gameState)
			return gameState, err
		},
		func() interface{} {
			return GameState{}
		}, func(value interface{}) interface{} {
			return update(value.(GameState))
		})
}

func (r *Repository) SubscribeGameState(ctx Context,
	gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (channel chan GameState, err error) {

	valid, err := r.CheckPlayerKey(ctx, gameId, playerId, playerKey)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidPlayerKeyErr
	}

	var gameState GameState
	subscription, err := redisUtil.SubscribeJson(ctx, r.RedisClient, GetGameStateKey(gameId), GetGameStatePubSubKey(gameId), &gameState)

	channel = make(chan GameState)
	go func() {
		defer close(channel)
		for {
			latest, ok := <-subscription
			if !ok {
				return
			}
			channel <- *latest.(*GameState)
		}
	}()
	return channel, nil
}
