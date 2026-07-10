package cache

import (
	"context"
	"source-base/pkg/dto"
	"time"
)

// Engine defines the standard interface for caching operations
type Engine interface {
	// Get retrieves a value by key.
	Get(ctx context.Context, key string) ([]byte, bool, error)

	// Set stores a value with an optional TTL.
	// value can be any serializable type.
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Delete removes a key from the cache.
	Delete(ctx context.Context, key string) error

	// BatchSet stores multiple values in a pipeline.
	// values is a map of key -> value.
	BatchSet(ctx context.Context, values map[string]any, ttl time.Duration) error

	// BatchDelete removes multiple keys from the cache.
	BatchDelete(ctx context.Context, keys []string) error

	// Close closes the connection to the cache server.
	Close() error
}

type SortedSetEngine interface {
	ZAdd(ctx context.Context, key string, members ...*dto.ZMember) error
	ZRemRangeByScore(ctx context.Context, key string, min, max string) error
	ZCount(ctx context.Context, key string, min, max string) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
}
