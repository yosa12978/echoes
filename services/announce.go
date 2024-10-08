package services

import (
	"context"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Get(ctx context.Context) (*types.Announce, error)
	Create(ctx context.Context, content string) error
	Delete(ctx context.Context) error
}

type announce struct {
	announceRepo repos.Announce
	logger       logging.Logger
}

func NewAnnounce(announceRepo repos.Announce, logger logging.Logger) Announce {
	return &announce{announceRepo: announceRepo, logger: logger}
}

func (s *announce) Get(ctx context.Context) (*types.Announce, error) {
	return s.announceRepo.Get(ctx)
}

func (s *announce) Create(ctx context.Context, content string) error {
	return s.announceRepo.Create(ctx, content)
}

func (s *announce) Delete(ctx context.Context) error {
	return s.announceRepo.Delete(ctx)
}
