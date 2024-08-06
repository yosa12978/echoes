package data

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/config"
)

var (
	redisOnce sync.Once
	rdb       *redis.Client
)

func Redis(ctx context.Context) *redis.Client {
	redisOnce.Do(func() {
		config := config.Get()
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.Db,
		})
		if err := rdb.Ping(ctx).Err(); err != nil {
			panic(err)
		}
	})
	return rdb
}

type redisPinger struct {
	rdb *redis.Client
}

func NewRedisPinger(ctx context.Context) Pinger {
	return &redisPinger{
		rdb: rdb,
	}
}

func (p *redisPinger) Ping(ctx context.Context) error {
	return p.rdb.Ping(ctx).Err()
}
