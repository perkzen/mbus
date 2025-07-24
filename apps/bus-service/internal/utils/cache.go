package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func TryGetFromCache[T any](ctx context.Context, client *redis.Client, key string) (*T, bool) {
	result, err := client.Get(ctx, key).Result()
	if err != nil || result == "" {
		return nil, false
	}

	var value T
	if err := json.Unmarshal([]byte(result), &value); err != nil {
		return nil, false
	}
	return &value, true
}

func SaveToCache(ctx context.Context, client *redis.Client, key string, value any, ttl time.Duration) {
	if data, err := json.Marshal(value); err == nil {
		_ = client.Set(ctx, key, data, ttl).Err()
	}
}

func WithCache[T any](ctx context.Context, cache *redis.Client, key string, ttl time.Duration, loader func() (T, error)) (T, error) {
	if cached, ok := TryGetFromCache[T](ctx, cache, key); ok {
		return *cached, nil
	}
	data, err := loader()
	if err != nil {
		var zero T
		return zero, err
	}
	SaveToCache(ctx, cache, key, data, ttl)
	return data, nil
}
