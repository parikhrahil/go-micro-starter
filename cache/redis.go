package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/parikhrahil/go-micro-starter/logger"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type Opts struct {
	Context   context.Context
	Log       logger.Logger
	RedisOpts *redis.Options
}

func New(opts *Opts) (CacheService, error) {
	ctx := opts.Context
	log := opts.Log
	redisOpts := opts.RedisOpts

	client := redis.NewClient(redisOpts)
	r := &RedisCache{client: client}

	if err := r.ping(ctx); err != nil {
		log.Error(fmt.Sprintf("Failed to connect to redis: %v", err))
		return nil, fmt.Errorf("Failed to connect to redis: %v", err)
	}
	log.Info(fmt.Sprintf("Redis connected at: %s", redisOpts.Addr))
	return r, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	if r == nil || r.client == nil {
		return fmt.Errorf("cannot close: redis is not initialized")
	}
	return r.client.Close()
}
