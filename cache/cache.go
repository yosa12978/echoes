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
	SortedSet
	Tx() (Pipeline, error)
}

type basic interface {
	Exists(ctx context.Context, keys ...string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Expires(ctx context.Context, key string, expiration time.Duration) (bool, error)
}

type String interface {
	basic
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error)
	SetXX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error)
}

type Hashmap interface {
	basic
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, value ...interface{}) (int64, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
}

type SortedSet interface {
	basic
	ZCard(ctx context.Context, key string) (int64, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZAdd(ctx context.Context, key string, members ...Member) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevrange(ctx context.Context, key string, start, stop int64) ([]string, error)
}

type Member struct {
	Score  float64
	Member string
}

type redisCache struct {
	rdb redis.Cmdable
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

func (c *redisCache) SetXX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error) {
	return c.rdb.SetXX(ctx, key, value, exp).Result()
}

func (c *redisCache) SetNX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, exp).Result()
}

func (c *redisCache) Del(ctx context.Context, keys ...string) (int64, error) {
	res, err := c.rdb.Del(ctx, keys...).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrNotFound
	}
	return res, err
}

func (c *redisCache) Expires(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.rdb.Expire(ctx, key, expiration).Result()
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
	res, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	return res, nil
}

func (c *redisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.rdb.Scan(ctx, cursor, match, count).Result()
}

func (c *redisCache) ZAdd(ctx context.Context, key string, members ...Member) (int64, error) {
	rmembers := make([]redis.Z, len(members))
	for k, v := range members {
		rmembers[k] = redis.Z{Score: v.Score, Member: v.Member}
	}
	return c.rdb.ZAdd(ctx, key, rmembers...).Result()
}

func (c *redisCache) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return c.rdb.ZScore(ctx, key, member).Result()
}

func (c *redisCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRange(ctx, key, start, stop).Result()
}

func (c *redisCache) ZRevrange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRevRange(ctx, key, start, stop).Result()
}

func (c *redisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Exists(ctx, keys...).Result()
}

func (c *redisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

func (c *redisCache) ZCard(ctx context.Context, key string) (int64, error) {
	return c.rdb.ZCard(ctx, key).Result()
}

func (c *redisCache) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.rdb.ZRem(ctx, key, members...).Result()
}

func (c *redisCache) Decr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Decr(ctx, key).Result()
}

func (c *redisCache) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return c.rdb.MGet(ctx, keys...).Result()
}

func (c *redisCache) Tx() (Pipeline, error) {
	return &redisTransaction{
		pipeline: c.rdb.TxPipeline(),
	}, nil
}

type Pipeline interface {
	Append(ctx context.Context, f func(pipe Cache) error) error
	Exec(ctx context.Context) error
	Discard(ctx context.Context) error
}

type redisTransaction struct {
	pipeline redis.Pipeliner
}

func (t *redisTransaction) Append(ctx context.Context, f func(Cache) error) error {
	return f(&redisCache{rdb: t.pipeline})
}

func (t *redisTransaction) Exec(ctx context.Context) error {
	_, err := t.pipeline.Exec(ctx)
	return err
}

func (t *redisTransaction) Discard(ctx context.Context) error {
	t.pipeline.Discard()
	return nil
}
