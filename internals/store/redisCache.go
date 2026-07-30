package store

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(redisClient *redis.Client) Cache {
	return &RedisCache{
		rdb: redisClient,
	}
}

func (rc *RedisCache) get(ctx context.Context, key string) (string, error) {
	val, err := rc.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rc *RedisCache) set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return rc.rdb.Set(ctx, key, value, ttl).Err()
}
