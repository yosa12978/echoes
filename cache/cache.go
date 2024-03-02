package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/data"
)

var ErrNotFound = errors.New("key doesn't exist")

type Cache interface {
	String
	Hashmap
}

type del interface {
	Del(ctx context.Context, keys ...string) (int64, error)
}

type String interface {
	del
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error)
}

type Hashmap interface {
	del
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, value ...interface{}) (int64, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}

type SortedSet interface {
	ZAdd(ctx context.Context, key string, members ...Z) (int64, error)
	ZScore(ctx context.Context, key string, member string) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]Z, error)
	ZRevrange(ctx context.Context, key string, start, stop int64) ([]Z, error)
}

type Z struct {
	Score  float64
	Member string
}

type redisCache struct {
	rdb *redis.Client
}

func NewRedisCache(ctx context.Context) Cache {
	return &redisCache{rdb: data.Redis(ctx)}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return res, err
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error) {
	res, err := c.rdb.Set(ctx, key, value, exp).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return res, err
}

func (c *redisCache) Del(ctx context.Context, keys ...string) (int64, error) {
	res, err := c.rdb.Del(ctx, keys...).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrNotFound
	}
	return res, err
}

func (c *redisCache) HGet(ctx context.Context, key, field string) (string, error) {
	res, err := c.rdb.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return res, err
}

func (c *redisCache) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	res, err := c.rdb.HSet(ctx, key, values...).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrNotFound
	}
	return res, err
}

func (c *redisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.rdb.HGetAll(ctx, key).Result()
}
