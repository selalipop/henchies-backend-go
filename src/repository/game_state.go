package repository

import (
	"context"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/redisutil"
	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

var gameStateRedisTTL = 5 * time.Hour
var maxElapsedGameUpdateTime = 5 * time.Minute

// InitGameState initializes game state with given parameters
// Returns an error if the game has already been initialized
func (r *Repository) InitGameState(ctx context.Context, gameID models.GameID, startingPlayerCount int, imposterCount int) error {
	exists, err := r.checkIfGameExists(ctx, gameID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("game already initialized")
	}

	err = internalUpdateGameStateTx(ctx, r.RedisClient, gameID, func(gameState models.GameState) models.GameState {
		return models.GameState{
			ImposterCount: imposterCount,
			MaxPlayers:    startingPlayerCount,
			Phase:         models.WaitingForPlayers,
			Players:       []models.PlayerID{},
		}
	})
	if err != nil {
		return errors.Wrap(err, "failed to initialize game state")
	}
	return nil
}

func (r *Repository) checkIfGameExists(ctx context.Context, gameID models.GameID) (bool, error) {
	exists, err := r.RedisClient.Exists(ctx, redisKeyGameState(gameID)).Result()
	if err != nil {
		return false, errors.Wrap(err, "error checking if game already exists")
	}
	return exists > 0, nil
}

// AddPlayerToGame adds a player to an existing game
// Will update game state and player state to reflect current game
func (r *Repository) AddPlayerToGame(ctx context.Context, gameID models.GameID, playerID models.PlayerID) error {
	exists, err := r.checkIfGameExists(ctx, gameID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("attempt to add player to non-existent game")
	}

	err = internalUpdateGameStateTx(ctx, r.RedisClient, gameID, func(gameState models.GameState) models.GameState {
		if gameState.Players.Contains(playerID) {
			return gameState
		}
		gameState.Players = append(gameState.Players, playerID)
		return gameState
	})
	if err != nil {
		return errors.Wrap(err, "failed to add player to game")
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

// ClearGameState clears game state transactionally
func (r *Repository) ClearGameState(ctx context.Context, gameID models.GameID) error {
	operation := func() error {
		return internalUpdateGameStateTxPtr(ctx, r.RedisClient, gameID, func(gameState *models.GameState) *models.GameState {
			for _, playerID := range gameState.Players {
				go func(currentPlayer models.PlayerID) {
					err := r.ClearPlayerState(ctx, gameID, currentPlayer)
					if err != nil {
						logrus.Errorf("failed to clear player state while clearing game state %v", err)
					}
				}(playerID)
			}
			return nil
		})
	}

	txBackoff := backoff.NewExponentialBackOff()
	txBackoff.MaxElapsedTime = maxElapsedGameUpdateTime
	return backoff.Retry(operation, txBackoff)
}

// SubscribeGameState returns a channel that passes on updates to Game State
// Will immediately return current game state
// Returns UnauthorizedPlayer error if the player key is not authorized to subscribe to this game
func (r *Repository) SubscribeGameState(ctx context.Context, gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey) (finished chan struct{}, channel chan models.GameState, err error) {
	err = r.CheckPlayerKey(ctx, gameID, playerID, playerKey)
	if err != nil {
		return finished, nil, err
	}

	var gameState models.GameState
	finished, subscription, err := redisutil.SubscribeJSON(ctx, r.RedisClient, redisKeyGameState(gameID), redisPublishKeyGameState(gameID), &gameState)
	if err != nil {
		return finished, nil, err
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
	return finished, channel, nil
}
func internalUpdateGameStateTx(ctx context.Context, client *redis.Client, gameID models.GameID, update func(gameState models.GameState) models.GameState) error {
	return internalUpdateGameStateTxPtr(ctx, client, gameID, func(gameState *models.GameState) *models.GameState {
		newValue := update(*gameState)
		return &newValue
	})
}
func internalUpdateGameStateTxPtr(ctx context.Context, client *redis.Client, gameID models.GameID, update func(gameState *models.GameState) *models.GameState) error {
	stateKey := redisKeyGameState(gameID)
	publishKey := redisPublishKeyGameState(gameID)

	var gameState models.GameState

	return redisutil.UpdateKeyTransaction(ctx, client, stateKey, publishKey, gameStateRedisTTL, 0, &gameState,
		func(value interface{}) interface{} {
			return update(value.(*models.GameState))
		})
}
