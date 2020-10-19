package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/util"
	. "github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

type GameRepository struct {
	PlayerRepo *PlayerRepository
	RepositoryEnv
}

func GetGameStatePubSubKey(gameId GameId) string {
	return fmt.Sprintf("playerGamePubSubKey:%s", gameId)
}
func GetGameStateKey(gameId GameId) string {
	return fmt.Sprintf("playerGameKey:%s", gameId)
}

var GameStateTTL = 5 * time.Hour
var MaxElapsedTimeGameStateUpdate = 5 * time.Minute

func (r *GameRepository) InitGameState(gameId GameId, startingPlayerCount int, imposterCount int) error {
	exists, err :=
		r.RedisClient.Exists(r.Context, GetGameStateKey(gameId)).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	return UpdateGameStateTransaction(gameId, r.RedisClient, r.Context, func(gameState GameState) GameState {
		return GameState{
			ImposterCount: imposterCount,
			MaxPlayers:    startingPlayerCount,
			Phase:         WaitingForPlayers,
			Players:       []PlayerId{},
		}
	})
}
func (r *GameRepository) AddPlayerToGame(gameId GameId, playerId PlayerId) error {
	exists, err :=
		r.RedisClient.Exists(r.Context, GetGameStateKey(gameId)).Result()

	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	err = UpdateGameStateTransaction(gameId, r.RedisClient, r.Context, func(gameState GameState) GameState {
		gameState.Players = append(gameState.Players, playerId)
		return gameState
	})
	if err != nil {
		return err
	}

	return r.PlayerRepo.UpdatePlayerStateUnchecked(gameId, playerId, func(state *PlayerState) {
		state.CurrentGame = gameId
	})
}
func (r *GameRepository) UpdateGameState(gameId GameId, update func(gameState GameState) GameState) error {

	operation := func() error {
		return UpdateGameStateTransaction(gameId, r.RedisClient, r.Context, update)
	}

	backoff := NewExponentialBackOff()
	backoff.MaxElapsedTime = MaxElapsedTimeGameStateUpdate
	return Retry(operation, backoff)
}

func UpdateGameStateTransaction(gameId GameId, client *redis.Client, context context.Context, update func(gameState GameState) GameState) error {
	stateKey := GetGameStateKey(gameId)
	publishKey := GetGameStatePubSubKey(gameId)

	return client.Watch(context, func(tx *redis.Tx) error {
		gameStateJson, err := tx.Get(context, stateKey).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		var currentState GameState
		if err != redis.Nil {
			err = json.Unmarshal([]byte(gameStateJson), &currentState)
			if err != nil {
				return err
			}
		}
		newState := update(currentState)
		newStateSerialized, err := json.Marshal(newState)
		_, err = tx.TxPipelined(context, func(pipe redis.Pipeliner) error {
			pipe.Set(context, stateKey, newStateSerialized, GameStateTTL)
			pipe.Publish(context, publishKey, newStateSerialized)
			return nil
		})
		return err
	}, stateKey)
}

func (r *GameRepository) SubscribeGameState(gameId GameId, playerId PlayerId, playerKey PlayerGameKey) (channel chan GameState, err error) {

	valid, err := r.CheckPlayerKey(gameId, playerId, playerKey)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidPlayerKeyErr
	}

	listen := r.RedisClient.Subscribe(r.Context, GetGameStatePubSubKey(gameId)).Channel()

	send := make(chan GameState)

	go func() {
		defer close(send)
		var gameState GameState

		err := util.GetRedisJson(r.Context, r.RedisClient, GetGameStateKey(gameId), &gameState)
		if err != nil {
			logrus.Error("failed to unmarshal Game State from key", err)
			return
		}
		send <- gameState
		for {
			message, ok := <-listen
			if !ok {
				return
			}

			err := json.Unmarshal([]byte(message.Payload), &gameState)
			if err != nil {
				logrus.Error("failed to unmarshal Game State from pubsub", err)
				break
			}
			send <- gameState
		}
	}()

	return send, nil
}
