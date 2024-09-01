package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
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

	//Deprecated
	GetLinkById(ctx context.Context, id string) (*types.Link, error)
	//Deprecated
	AddLink(ctx context.Context, link types.Link) error
	//Deprecated
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

// Deprecated
func addLinkCmdable(ctx context.Context, link types.Link, rdb redis.Cmdable) error {
	linkMap := map[string]interface{}{
		"id":      link.Id,
		"name":    link.Name,
		"created": link.Created,
		"icon":    link.Icon,
		"place":   link.Place,
	}
	if err := rdb.HSet(ctx, "links:"+link.Id, linkMap).Err(); err != nil {
		return newInternalFailure(err)
	}
	if err := rdb.Del(ctx, "links").Err(); err != nil { // refresh links on main page
		if errors.Is(err, redis.Nil) {
			return newNotFound(err)
		}
		return newInternalFailure(err)
	}
	return nil
}

// Deprecated
func (l *linkCache) AddLink(ctx context.Context, link types.Link) error {
	pipe := l.rdb.Pipeline()
	if err := addLinkCmdable(ctx, link, pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (l *linkCache) Delete(ctx context.Context, id string) error {
	// pipe := l.rdb.Pipeline()
	// if err := pipe.Del(ctx, "links:"+id).Err(); err != nil {
	// 	if errors.Is(err, redis.Nil) {
	// 		return newNotFound(err)
	// 	}
	// 	return newInternalFailure(err)
	// }
	// if err := pipe.Del(ctx, "links").Err(); err != nil { // refresh links on main page
	// 	if errors.Is(err, redis.Nil) {
	// 		return newNotFound(err)
	// 	}
	// 	return newInternalFailure(err)
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

// Deprecated
func (l *linkCache) GetLinkById(ctx context.Context, id string) (*types.Link, error) {
	linkMap, err := l.rdb.HGetAll(ctx, "links:"+id).Result()
	if err != nil {
		return nil, newInternalFailure(err)
	}
	if len(linkMap) == 0 {
		return nil, newNotFound(err)
	}

	place, _ := strconv.Atoi(linkMap["place"])
	link := types.Link{
		Id:      linkMap["id"],
		Name:    linkMap["name"],
		Created: linkMap["created"],
		Icon:    linkMap["icon"],
		Place:   place,
	}
	return &link, nil
}

func (l *linkCache) AddLinks(ctx context.Context, links ...types.Link) error {
	linksJson, _ := json.Marshal(links)
	if err := l.rdb.Set(ctx, "links", linksJson, 90*time.Second).Err(); err != nil {
		return newInternalFailure(err)
	}
	return nil
}

func (l *linkCache) GetLinks(ctx context.Context) ([]types.Link, error) {
	linksJson, err := l.rdb.Get(ctx, "links").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, newNotFound(err)
		}
		return nil, newInternalFailure(err)
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
			return newNotFound(err)
		}
		return newInternalFailure(err)
	}
	return nil
}
