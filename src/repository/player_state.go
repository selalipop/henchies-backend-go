package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/SelaliAdobor/henchies-backend-go/src/redisutil"
	"github.com/go-redis/redis/v8"
	"time"
)

const playerIDTTL = 7200 * time.Second

// SubscribePlayerState returns a channel that passes on updates to Player State
// Will immediately return current player state state
// Returns UnauthorizedPlayer error if the player key is not authorized to subscribe to this player
func (r *Repository) SubscribePlayerState(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey) (closeChannel chan struct{}, channel chan models.PlayerState, err error) {
	err = r.CheckPlayerKey(ctx, gameID, playerID, playerKey)
	if err != nil {
		return nil, nil, err
	}

	stateKey := redisKeyPlayerState(gameID, playerID)
	publishKey := redisKeyPlayerStatePublish(gameID, playerID)

	var playerState models.PlayerState

	finished, subscription, err := redisutil.SubscribeJSON(ctx, r.RedisClient, stateKey, publishKey, &playerState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to player state in redis %w", err)
	}

	channel = make(chan models.PlayerState)
	go func() {
		defer close(channel)
		for {
			latest, ok := <-subscription
			if !ok {
				return
			}
			channel <- *latest.(*models.PlayerState)
		}
	}()
	return finished, channel, nil
}

// GetPlayerStateChecked retrieves current player state, checks for a valid player key
func (r Repository) GetPlayerStateChecked(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey) (state models.PlayerState, err error) {
	return r.internalGetPlayerState(ctx, gameID, playerID, playerKey, true)
}

// GetPlayerStateUnchecked retrieves current player state without checking ID, for internal use only, do not allow players to access unchecked player state
func (r Repository) GetPlayerStateUnchecked(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID) (state models.PlayerState, err error) {
	return r.internalGetPlayerState(ctx, gameID, playerID, models.PlayerGameKey{}, false)
}

// UpdatePlayerStateChecked updates current player state while checking ID
// Returns UnauthorizedPlayerError if the key passed did not match the given player ID
func (r Repository) UpdatePlayerStateChecked(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey, update func(state models.PlayerState) models.PlayerState) error {
	return r.internalPlayerStateUpdate(ctx, gameID, playerID, playerKey, true, r, update)
}

// UpdatePlayerStateUnchecked updates current player state without checking ID
// For server use only, do not allow players to access unchecked player state
func (r Repository) UpdatePlayerStateUnchecked(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, update func(state models.PlayerState) models.PlayerState) error {
	return r.internalPlayerStateUpdate(ctx, gameID, playerID, models.PlayerGameKey{}, false, r, update)
}

// ClearPlayerState clears player's data including current game
// For server use only, do not allow players to clear player state
func (r Repository) ClearPlayerState(ctx context.Context, gameID models.GameID, playerID models.PlayerID) error {
	return r.UpdatePlayerStateUnchecked(ctx, gameID, playerID, func(state models.PlayerState) models.PlayerState {
		if state.CurrentGame != gameID {
			return state
		}
		return models.PlayerState{
			CurrentGame: "",
			GameKey:     models.PlayerGameKey{},
			IsImposter:  false,
		}
	})
}

// CheckPlayerKey compares the given key to the one stored for the player
// If there is no error and the 'valid' is false, this was an attempt to access player state with the wrong key
func (r Repository) CheckPlayerKey(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey) error {
	redisKey := redisKeyPlayerGameKey(gameID, playerID)
	storedKeyJSON, err := r.RedisClient.Get(ctx, redisKey).Result()

	if err != nil {
		return fmt.Errorf("failed to get player game key from redis (%v): %w", redisKey, err)
	}
	var storedKey models.PlayerGameKey
	err = json.Unmarshal([]byte(storedKeyJSON), &storedKey)
	if err != nil {
		return fmt.Errorf("failed to unmarshal player game key from redis (%v): %w", redisKey, err)
	}

	if storedKey.Key != playerKey.Key {
		return UnauthorizedPlayer
	}
	return nil
}

func (r Repository) internalGetPlayerState(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey, shouldCheck bool) (state models.PlayerState, err error) {
	if shouldCheck {
		err := r.CheckPlayerKey(ctx, gameID, playerID, playerKey)
		if err != nil {
			return state, err
		}
	}

	var keyName = redisKeyPlayerState(gameID, playerID)

	err = redisutil.GetRedisJSON(ctx, r.RedisClient, keyName, &state)

	if err != nil {
		return state, fmt.Errorf("failed to get player state from redis: %w", err)
	}
	return
}

func (r Repository) internalPlayerStateUpdate(ctx context.Context,
	gameID models.GameID, playerID models.PlayerID, playerKey models.PlayerGameKey, shouldCheck bool, env Repository, update func(state models.PlayerState) models.PlayerState) error {
	if shouldCheck {
		err := env.CheckPlayerKey(ctx, gameID, playerID, playerKey)
		if err != nil {
			return err
		}
	}

	return internalPlayerStateUpdateTransaction(ctx, env.RedisClient, gameID, playerID, update)
}

func internalPlayerStateUpdateTransaction(ctx context.Context,
	client *redis.Client, gameID models.GameID, playerID models.PlayerID, update func(state models.PlayerState) models.PlayerState) error {
	stateKey := redisKeyPlayerState(gameID, playerID)
	publishKey := redisKeyPlayerStatePublish(gameID, playerID)

	var playerState models.PlayerState

	return redisutil.UpdateKeyTransaction(ctx, client, stateKey, publishKey, gameStateRedisTTL, 0, &playerState,
		func(value interface{}) interface{} {
			return update(*value.(*models.PlayerState))
		})
}
