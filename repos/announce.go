package repos

import (
	"context"
	"time"

	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Get(ctx context.Context) (*types.Announce, error)
	Create(ctx context.Context, content string) error
	Delete(ctx context.Context) error
}

type announce struct {
	storage *types.Announce
}

func NewAnnounce() Announce {
	repo := new(announce)
	repo.storage = nil
	return repo
}

func (repo *announce) Create(ctx context.Context, content string) error {
	repo.storage = new(types.Announce)
	repo.storage.Content = content
	repo.storage.Date = time.Now().Format(time.RFC3339)
	return nil
}

func (repo *announce) Delete(ctx context.Context) error {
	repo.storage = nil
	return nil
}

func (repo *announce) Get(ctx context.Context) (*types.Announce, error) {
	return repo.storage, nil
}

type announceCacheAdapter struct {
	cache cache.Announce
}

func NewAnnounceCacheAdapter(cache cache.Announce) Announce {
	return &announceCacheAdapter{
		cache: cache,
	}
}

func (a *announceCacheAdapter) Create(ctx context.Context, content string) error {
	return a.cache.Create(ctx, content)
}

func (a *announceCacheAdapter) Delete(ctx context.Context) error {
	return a.cache.Delete(ctx)
}

func (a *announceCacheAdapter) Get(ctx context.Context) (*types.Announce, error) {
	return a.cache.Get(ctx)
}
