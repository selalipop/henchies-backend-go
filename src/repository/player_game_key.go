package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SelaliAdobor/henchies-backend-go/src/models"
	"github.com/google/uuid"
)

// GetPlayerGameKey returns the game key for a given player for a specific game
// The IP address must match the first IP address used to access this key
// If the IP used to access to endpoint changes, an error is returned
// This prevents other players from using a player's game key
func (r *Repository) GetPlayerGameKey(ctx context.Context, gameID models.GameID, playerID models.PlayerID, ipAddress string) (playerKey models.PlayerGameKey, err error) {
	var keyName = redisKeyPlayerGameKey(gameID, playerID)

	newPlayerID := uuid.New()
	newKey := models.PlayerGameKey{Key: newPlayerID.String(), OwnerIP: ipAddress}

	serializedKey, err := json.Marshal(newKey)
	if err != nil {
		return playerKey, fmt.Errorf("failed to marshal new player key %w", err)
	}

	_, err = r.RedisClient.SetNX(ctx, keyName, serializedKey, playerIDTTL).Result()
	if err != nil {
		return playerKey, fmt.Errorf("failed to setnx player game key from redis %w", err)
	}

	result, err := r.RedisClient.Get(ctx, keyName).Result()
	if err != nil {
		return playerKey, fmt.Errorf("failed to get player game key from redis %w", err)
	}

	var currentKey models.PlayerGameKey

	err = json.Unmarshal([]byte(result), &currentKey)
	if err != nil {
		return playerKey, fmt.Errorf("failed to get player game key from redis %w", err)
	}

	if currentKey == newKey {
		err := r.UpdatePlayerStateUnchecked(ctx, gameID, playerID, func(state models.PlayerState) models.PlayerState {
			state.GameKey = currentKey
			return state
		})
		if err != nil {
			return playerKey, fmt.Errorf("failed to update player state with new key %w", err)
		}
	}

	if currentKey.OwnerIP != ipAddress {
		return playerKey, errors.New("ip address mismatch retrieving key")
	}
	return currentKey, nil
}
