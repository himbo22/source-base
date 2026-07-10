package cache

import (
	"context"
	"encoding/json"
	"fmt"
)

func Load[T any](
	ctx context.Context,
	cache Engine,
	key string,
	dest *T,
) (bool, error) {
	byteData, exists, err := cache.Get(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	// cache miss
	if !exists {
		return false, nil
	}

	if err := json.Unmarshal(byteData, dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	return true, nil
}

//  usage
//	var user User
//
//	found, err := cache.Load(ctx, c, key, &user)
//	if err != nil {
//	return nil, err
//	}
//
//	if found {
//	return &user, nil
//	}
//
//	// Cache miss -> query DB
//	user, err = repo.FindByID(ctx, id)
