package store

import (
	"context"
	"time"
)
type Cache interface {
	get(ctx context.Context, key string) (string, error)
	set(ctx context.Context, key string, value string, ttl time.Duration) error
}