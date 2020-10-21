package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/redisutil"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v8"
	"time"
)

var gameStateRedisTTL = 5 * time.Hour
var maxElapsedGameUpdateTime = 5 * time.Minute

// InitGameState initializes game state with given parameters
// Returns an error if the game has already been initialized
func (r *Repository) InitGameState(ctx context.Context, gameID models.GameID, startingPlayerCount int, imposterCount int) error {
	exists, err :=
		r.RedisClient.Exists(ctx, redisKeyGameState(gameID)).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	return internalUpdateGameStateTx(ctx, r.RedisClient, gameID, func(gameState models.GameState) models.GameState {
		return models.GameState{
			ImposterCount: imposterCount,
			MaxPlayers:    startingPlayerCount,
			Phase:         models.WaitingForPlayers,
			Players:       []models.PlayerID{},
		}
	})
}

// AddPlayerToGame adds a player to an existing game
// Will update game state and player state to reflect current game
func (r *Repository) AddPlayerToGame(ctx context.Context, gameID models.GameID, playerID models.PlayerID) error {
	exists, err :=
		r.RedisClient.Exists(ctx, redisKeyGameState(gameID)).Result()

	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("game already initialized")
	}

	err = internalUpdateGameStateTx(ctx, r.RedisClient, gameID, func(gameState models.GameState) models.GameState {
		gameState.Players = append(gameState.Players, playerID)
		return gameState
	})
	if err != nil {
		return err
	}

	return r.UpdatePlayerStateUnchecked(ctx, gameID, playerID, func(state models.PlayerState) models.PlayerState {
		state.CurrentGame = gameID
		return state
	})
}

// UpdateGameState updates game state transactionally
// Warning: update function may be called multiple times if the game state is modified while it is running
// Update function should contain preconditions to verify current state before modifications
func (r *Repository) UpdateGameState(ctx context.Context, gameID models.GameID, update func(gameState models.GameState) models.GameState) error {
	operation := func() error {
		return internalUpdateGameStateTx(ctx, r.RedisClient, gameID, update)
	}

	txBackoff := backoff.NewExponentialBackOff()
	txBackoff.MaxElapsedTime = maxElapsedGameUpdateTime
	return backoff.Retry(operation, txBackoff)
}

// SubscribeGameState returns a channel that passes on updates to Game State
// Will immediately return current game state
// Returns UnauthorizedPlayer error if the player key is not authorized to subscribe to this game
func (r *Repository) SubscribeGameState(
	ctx context.Context,
	gameID models.GameID,
	playerID models.PlayerID,
	playerKey models.PlayerGameKey,
) (channel chan models.GameState, err error) {
	err = r.CheckPlayerKey(ctx, gameID, playerID, playerKey)
	if err != nil {
		return nil, err
	}

	var gameState models.GameState
	subscription, err := redisutil.SubscribeJSON(ctx, r.RedisClient, redisKeyGameState(gameID), redisPublishKeyGameState(gameID), &gameState)
	if err != nil {
		return nil, err
	}
	channel = make(chan models.GameState)
	go func() {
		defer close(channel)
		for {
			latest, ok := <-subscription
			if !ok {
				return
			}
			channel <- *latest.(*models.GameState)
		}
	}()
	return channel, nil
}

func internalUpdateGameStateTx(ctx context.Context, client *redis.Client, gameID models.GameID, update func(gameState models.GameState) models.GameState) error {
	stateKey := redisKeyGameState(gameID)
	publishKey := redisPublishKeyGameState(gameID)

	var gameState models.GameState

	return redisutil.UpdateKeyTransaction(ctx, client, stateKey, publishKey, gameStateRedisTTL, 0, &gameState,
		func(value interface{}) interface{} {
			return update(*value.(*models.GameState))
		})
}

func redisPublishKeyGameState(gameID models.GameID) string {
	return fmt.Sprintf("playerGamePubSubKey:%s", gameID)
}

func redisKeyGameState(gameID models.GameID) string {
	return fmt.Sprintf("playerGameKey:%s", gameID)
}
