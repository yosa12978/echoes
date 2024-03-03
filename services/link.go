package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	GetLinks(ctx context.Context) ([]types.Link, error)
	CreateLink(ctx context.Context, name, url string) (*types.Link, error)
	DeleteLink(ctx context.Context, id string) (*types.Link, error)
	Seed(ctx context.Context) error
}

type link struct {
	linkRepo repos.Link
	cache    cache.Hashmap
}

func NewLink(linkRepo repos.Link, cache cache.Hashmap) Link {
	return &link{linkRepo: linkRepo, cache: cache}
}

func (s *link) GetLinks(ctx context.Context) ([]types.Link, error) {
	return s.linkRepo.FindAll(ctx)
}

func (s *link) CreateLink(ctx context.Context, name, addr string) (*types.Link, error) {
	nameTrim := strings.TrimSpace(name)
	addrTrim := strings.TrimSpace(addr)
	if nameTrim == "" || addrTrim == "" {
		return nil, errors.New("link's name or url can't be an empty string")
	}
	_, err := url.ParseRequestURI(addrTrim)
	if err != nil {
		return nil, err
	}
	link := types.Link{
		Id:      uuid.NewString(),
		Name:    nameTrim,
		URL:     addr,
		Created: time.Now().Format(time.RFC3339),
	}
	return s.linkRepo.Create(ctx, link)
}

func (s *link) DeleteLink(ctx context.Context, id string) (*types.Link, error) {
	return s.linkRepo.Delete(ctx, id)
}

func (s *link) Seed(ctx context.Context) error {
	_, err := s.linkRepo.Create(ctx, types.Link{
		Id:      "09741221-7ea7-4106-ac19-8d2c2c90afbc",
		Name:    "reddit",
		URL:     "https://reddit.com",
		Created: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	_, err = s.linkRepo.Create(ctx, types.Link{
		Id:      "c46428bd-a807-4042-812b-f3b56f047732",
		Name:    "my github",
		URL:     "https://github.com/yosa12978",
		Created: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	_, err = s.linkRepo.Create(ctx, types.Link{
		Id:      "60a9f6e8-8fda-480a-832a-3e3a07ae8890",
		Name:    "wow forum (icy veins)",
		URL:     "https://www.icy-veins.com/",
		Created: time.Now().Format(time.RFC3339),
	})
	return err
}
