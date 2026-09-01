package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/himbo22/source-base/pkg/common/cache"
	"github.com/himbo22/source-base/pkg/dto"
	"github.com/himbo22/source-base/pkg/settings"
	"github.com/himbo22/source-base/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultPoolSize        = 10
	defaultMinIdleConns    = 5
	defaultPoolTimeout     = 5
	defaultDialTimeout     = 5
	defaultReadTimeout     = 3
	defaultWriteTimeout    = 3
	defaultMaxRetries      = 3
	defaultMinRetryBackoff = 300 // millis
	defaultMaxRetryBackoff = 500 // millis
)

type Engine struct {
	Client redis.UniversalClient
	Config *settings.Redis
}

var _ cache.Engine = (*Engine)(nil)
var _ cache.SortedSetEngine = (*Engine)(nil)

func (r *Engine) Connect() error {
	r.setDefaultConfig()

	r.Client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           r.Config.Addrs,
		MasterName:      r.Config.MasterName,
		Password:        r.Config.Password,
		DB:              r.Config.Database,
		PoolSize:        r.Config.PoolSize,
		MinIdleConns:    r.Config.MinIdleConns,
		MaxRetries:      r.Config.MaxRetries,
		DialTimeout:     utils.ToDuration(r.Config.DialTimeout),
		ReadTimeout:     utils.ToDuration(r.Config.ReadTimeout),
		WriteTimeout:    utils.ToDuration(r.Config.WriteTimeout),
		PoolTimeout:     utils.ToDuration(r.Config.PoolTimeout),
		MinRetryBackoff: utils.ToDurationMs(r.Config.MinRetryBackoff),
		MaxRetryBackoff: utils.ToDurationMs(r.Config.MaxRetryBackoff),
	})

	// ping
	ctx, cancel := context.WithTimeout(context.Background(), utils.ToDuration(r.Config.DialTimeout))
	defer cancel()
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// Get retrieves a value by key.
// found=false means key does not exist (not an error condition).
func (r *Engine) Get(ctx context.Context, key string) ([]byte, bool, error) {
	value, err := r.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("redis get key %q: %w", key, err)
	}
	return value, true, nil
}

// Exists checks if one or more keys exist in Redis. Returns the count of existing keys.
func (r *Engine) Exists(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	count, err := r.Client.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis exists keys %v: %w", keys, err)
	}
	return count, nil
}

// Set stores a value with the given TTL. ttl=0 means no expiration.
func (r *Engine) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value for key %q: %w", key, err)
	}
	if err := r.Client.Set(ctx, key, byteValue, ttl).Err(); err != nil {
		return fmt.Errorf("redis set key %q: %w", key, err)
	}
	return nil
}

// Delete removes a single key. Deleting a non-existent key is not an error.
func (r *Engine) Delete(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete key %q: %w", key, err)
	}
	return nil
}

// BatchSet stores multiple key-value pairs sharing the same TTL, in one pipeline round-trip.
func (r *Engine) BatchSet(ctx context.Context, values map[string]any, ttl time.Duration) error {
	if len(values) == 0 {
		return nil
	}

	pipe := r.Client.Pipeline()
	for key, value := range values {
		byteValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshal value for key %q: %w", key, err)
		}
		pipe.Set(ctx, key, byteValue, ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis batch set: %w", err)
	}
	return nil
}

// BatchDelete removes multiple keys in a single round-trip.
func (r *Engine) BatchDelete(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := r.Client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis batch delete: %w", err)
	}
	return nil
}

// ZAdd adds one or more members with scores to a sorted set.
func (r *Engine) ZAdd(ctx context.Context, key string, members ...*dto.ZMember) error {
	if len(members) == 0 {
		return nil
	}

	zMembers := make([]redis.Z, 0, len(members))
	for _, m := range members {
		zMembers = append(zMembers, redis.Z{
			Score:  m.Score,
			Member: m.Member,
		})
	}

	if err := r.Client.ZAdd(ctx, key, zMembers...).Err(); err != nil {
		return fmt.Errorf("redis zadd key %q: %w", key, err)
	}
	return nil
}

// ZRemRangeByScore removes members within the given score range [min, max].
// min/max accept redis range syntax, e.g. "-inf", "+inf", "(1", "5".
func (r *Engine) ZRemRangeByScore(ctx context.Context, key string, min, max string) error {
	if err := r.Client.ZRemRangeByScore(ctx, key, min, max).Err(); err != nil {
		return fmt.Errorf("redis zremrangebyscore key %q: %w", key, err)
	}
	return nil
}

// ZCount counts members within the given score range [min, max].
func (r *Engine) ZCount(ctx context.Context, key string, min, max string) (int64, error) {
	count, err := r.Client.ZCount(ctx, key, min, max).Result()
	if err != nil {
		return 0, fmt.Errorf("redis zcount key %q: %w", key, err)
	}
	return count, nil
}

// ZRange returns members within the given rank range [start, stop] (inclusive, 0-indexed).
func (r *Engine) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	members, err := r.Client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("redis zrange key %q: %w", key, err)
	}
	return members, nil
}

// Close releases the underlying connection pool. Safe to call once during shutdown.
func (r *Engine) Close() error {
	if r.Client == nil {
		return nil
	}
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("redis close: %w", err)
	}
	return nil
}

func (r *Engine) setDefaultConfig() {
	if r.Config.PoolSize == 0 {
		r.Config.PoolSize = defaultPoolSize
	}
	if r.Config.MinIdleConns == 0 {
		r.Config.MinIdleConns = defaultMinIdleConns
	}
	if r.Config.PoolTimeout == 0 {
		r.Config.PoolTimeout = defaultPoolTimeout
	}
	if r.Config.DialTimeout == 0 {
		r.Config.DialTimeout = defaultDialTimeout
	}
	if r.Config.ReadTimeout == 0 {
		r.Config.ReadTimeout = defaultReadTimeout
	}
	if r.Config.WriteTimeout == 0 {
		r.Config.WriteTimeout = defaultWriteTimeout
	}
	if r.Config.MaxRetries == 0 {
		r.Config.MaxRetries = defaultMaxRetries
	}
	if r.Config.MinRetryBackoff == 0 {
		r.Config.MinRetryBackoff = defaultMinRetryBackoff
	}
	if r.Config.MaxRetryBackoff == 0 {
		r.Config.MaxRetryBackoff = defaultMaxRetryBackoff
	}
}
