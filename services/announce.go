package services

import (
	"context"
	"errors"
	"strings"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Announce interface {
	Get(ctx context.Context) (*types.Announce, error)
	Create(ctx context.Context, content string) (*types.Announce, error)
	Delete(ctx context.Context) (*types.Announce, error)
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

func (s *announce) Create(ctx context.Context, content string) (*types.Announce, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("announce can't be empty")
	}
	return s.announceRepo.Create(ctx, strings.ReplaceAll(content, "\n", "<br>"))
}

func (s *announce) Delete(ctx context.Context) (*types.Announce, error) {
	return s.announceRepo.Delete(ctx)
}
