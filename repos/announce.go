package repos

import (
	"context"
	"time"

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
