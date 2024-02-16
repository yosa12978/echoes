package repos

import (
	"time"

	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Create(content string) types.Announce
	Delete()
	Get() *types.Announce
}

type announce struct {
	storage *types.Announce
}

func NewAnnounce() Announce {
	repo := new(announce)
	repo.storage = nil
	return repo
}

func (repo *announce) Create(content string) types.Announce {
	repo.storage = new(types.Announce)
	repo.storage.Content = content
	repo.storage.Date = time.Now().Format(time.RFC3339)
	return *repo.storage
}

func (repo *announce) Delete() {
	repo.storage = nil
}

func (repo *announce) Get() *types.Announce {
	return repo.storage
}
