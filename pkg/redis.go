package pkg

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	RedisClient *redis.Client
	ctx         context.Context
}

func NewCacheCache(ctx context.Context, address, password string, poolSize, minIdleConns, db int) (*Cache, error) {
	if poolSize < 1 {
		poolSize = 1000
	}
	if minIdleConns < 1 {
		minIdleConns = 10
	}
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		DB:           db,
		MinIdleConns: minIdleConns,
		PoolSize:     poolSize,
		Password:     password,
	})

	pong, err := client.Ping(ctx).Result()
	if pong != "PONG" || err != nil {
		return nil, err
	}

	return &Cache{
		RedisClient: client,
		ctx:         ctx,
	}, nil
}

func (s *Cache) Set(key string, value interface{}, expire time.Duration) (string, error) {
	return s.RedisClient.Set(s.ctx, key, value, expire).Result()
}

func (s *Cache) Get(key string) (string, error) {
	return s.RedisClient.Get(s.ctx, key).Result()
}

func (s *Cache) Del(keys ...string) (int64, error) {
	return s.RedisClient.Del(s.ctx, keys...).Result()
}

func (s *Cache) Incr(key string) (int64, error) {
	return s.RedisClient.Incr(s.ctx, key).Result()
}

func (s *Cache) Decr(key string) (int64, error) {
	return s.RedisClient.Decr(s.ctx, key).Result()
}

func (s *Cache) Exists(keys ...string) (int64, error) {
	return s.RedisClient.Exists(s.ctx, keys...).Result()
}
