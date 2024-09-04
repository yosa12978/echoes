package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Get(ctx context.Context) (*types.Announce, error)

	Create(ctx context.Context, content string) error
	Delete(ctx context.Context) error
}

type announceRedis struct {
	rdb *redis.Client
}

func NewAnnounceRedis(rdb *redis.Client) Announce {
	return &announceRedis{
		rdb: rdb,
	}
}

func (a *announceRedis) Get(ctx context.Context) (*types.Announce, error) {
	res, err := a.rdb.HGetAll(ctx, "announce").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, types.ErrNotFound
		}
		return nil, types.NewErrInternalFailure(err)
	}
	if len(res) == 0 {
		return nil, nil
	}
	announce := types.Announce{
		Content: res["content"],
		Date:    res["date"],
	}
	return &announce, nil
}

func (a *announceRedis) Create(ctx context.Context, content string) error {
	announce := map[string]string{
		"content": content,
		"date":    time.Now().Format(time.RFC3339),
	}
	_, err := a.rdb.HSet(ctx, "announce", announce).Result()
	if err != nil {
		return types.NewErrInternalFailure(err)
	}
	return nil
}

func (a *announceRedis) Delete(ctx context.Context) error {
	_, err := a.rdb.Del(ctx, "announce").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return types.ErrNotFound
		}
		return types.NewErrInternalFailure(err)
	}
	return err
}
