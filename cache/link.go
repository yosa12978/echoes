package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	GetLinks(ctx context.Context) ([]types.Link, error)
	AddLinks(ctx context.Context, links ...types.Link) error
	Flush(ctx context.Context) error
	Delete(ctx context.Context, id string) error
	GetLinkById(ctx context.Context, id string) (*types.Link, error)
	AddLink(ctx context.Context, link types.Link) error

	//Unimplemented
	Update(ctx context.Context, id string, link types.Link) error
}

type linkCache struct {
	rdb    *redis.Client
	logger logging.Logger
}

func NewLinkRedis(rdb *redis.Client, logger logging.Logger) Link {
	return &linkCache{
		rdb:    rdb,
		logger: logger,
	}
}

func (l *linkCache) AddLink(ctx context.Context, link types.Link) error {
	linkJson, _ := json.Marshal(link)
	if err := l.rdb.Set(ctx,
		"links:"+link.Id,
		linkJson,
		100*time.Second).Err(); err != nil {
		return types.NewErrInternalFailure(err)
	}
	return nil
}

func (l *linkCache) Delete(ctx context.Context, id string) error {
	// pipe := l.rdb.Pipeline()
	// if err := pipe.Del(ctx, "links:"+id).Err(); err != nil {
	// 	if errors.Is(err, redis.Nil) {
	// 		return types.NewErrNotFound(err)
	// 	}
	// 	return types.NewErrInternalFailure(err)
	// }
	// if err := pipe.Del(ctx, "links").Err(); err != nil { // refresh links on main page
	// 	if errors.Is(err, redis.Nil) {
	// 		return types.NewErrNotFound(err)
	// 	}
	// 	return types.NewErrInternalFailure(err)
	// }
	// _, err := pipe.Exec(ctx)
	links, err := l.GetLinks(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < len(links); i++ {
		if links[i].Id == id {
			links = append(links[:i], links[i+1:]...)
			break
		}
	}
	return l.AddLinks(ctx, links...)
}

func (l *linkCache) GetLinkById(ctx context.Context, id string) (*types.Link, error) {
	linkJson, err := l.rdb.Get(ctx, "links:"+id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, types.NewErrNotFound(err)
		}
		return nil, types.NewErrInternalFailure(err)
	}
	var link types.Link
	json.Unmarshal([]byte(linkJson), &link)
	return &link, nil
}

func (l *linkCache) AddLinks(ctx context.Context, links ...types.Link) error {
	linksJson, _ := json.Marshal(links)
	if err := l.rdb.Set(ctx, "links", linksJson, 90*time.Second).Err(); err != nil {
		return types.NewErrInternalFailure(err)
	}
	return nil
}

func (l *linkCache) GetLinks(ctx context.Context) ([]types.Link, error) {
	linksJson, err := l.rdb.Get(ctx, "links").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, types.NewErrNotFound(err)
		}
		return nil, types.NewErrInternalFailure(err)
	}

	var links []types.Link
	err = json.Unmarshal([]byte(linksJson), &links)
	return links, err
}

func (l *linkCache) Update(ctx context.Context, id string, link types.Link) error {
	panic("unimplemented")
}

func (l *linkCache) Flush(ctx context.Context) error {
	if err := l.rdb.Del(ctx, "links").Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return types.NewErrNotFound(err)
		}
		return types.NewErrInternalFailure(err)
	}
	return nil
}
