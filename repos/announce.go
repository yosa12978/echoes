package repos

import (
	"context"
	"errors"
	"time"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Get(ctx context.Context) (*types.Announce, error)
	Create(ctx context.Context, content string) (*types.Announce, error)
	Delete(ctx context.Context) (*types.Announce, error)
}

type announce struct {
	storage *types.Announce
}

func NewAnnounce() Announce {
	repo := new(announce)
	repo.storage = nil
	return repo
}

func (repo *announce) Create(ctx context.Context, content string) (*types.Announce, error) {
	repo.storage = new(types.Announce)
	repo.storage.Content = content
	repo.storage.Date = time.Now().Format(time.RFC3339)
	return repo.storage, nil
}

func (repo *announce) Delete(ctx context.Context) (*types.Announce, error) {
	repo.storage = nil
	return repo.storage, nil
}

func (repo *announce) Get(ctx context.Context) (*types.Announce, error) {
	return repo.storage, nil
}

type announceCache struct {
	cache cache.Hashmap
}

func NewAnnounceCache(cache cache.Hashmap) Announce {
	return &announceCache{
		cache: cache,
	}
}

func (a *announceCache) Get(ctx context.Context) (*types.Announce, error) {
	res, err := a.cache.HGetAll(ctx, "announce")
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	announce := types.Announce{
		Content: res["content"],
		Date:    res["date"],
	}
	return &announce, nil
}

func (a *announceCache) Create(ctx context.Context, content string) (*types.Announce, error) {
	announce := map[string]string{
		"content": content,
		"date":    time.Now().Format(time.RFC3339),
	}
	_, err := a.cache.HSet(ctx, "announce", announce)
	res := types.Announce{Content: content, Date: time.Now().Format(time.RFC3339)}
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &res, nil
}

func (a *announceCache) Delete(ctx context.Context) (*types.Announce, error) {
	announce, err := a.Get(ctx)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	_, err = a.cache.Del(ctx, "announce")
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return announce, err
}
