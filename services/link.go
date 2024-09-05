package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	GetLinks(ctx context.Context) ([]types.Link, error)
	GetLinkById(ctx context.Context, id string) (*types.Link, error)
	CreateLink(ctx context.Context, name, url, icon string, place int) (*types.Link, error)
	DeleteLink(ctx context.Context, id string) (*types.Link, error)
	Seed(ctx context.Context) error
}

type link struct {
	linkRepo repos.Link
	cache    cache.Link
	logger   logging.Logger
}

func NewLink(linkRepo repos.Link, cache cache.Link, logger logging.Logger) Link {
	return &link{linkRepo: linkRepo, cache: cache, logger: logger}
}

func (s *link) GetLinks(ctx context.Context) ([]types.Link, error) {
	links, err := s.cache.GetLinks(ctx)
	if err != nil {
		if errors.Is(err, types.ErrInternalFailure) {
			s.logger.Error(err.Error())
		}
	}
	if links != nil {
		return links, nil
	}

	links, err = s.linkRepo.FindAll(ctx)
	if err != nil {
		return links, err
	}

	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.cache.AddLinks(timeout, links...); err != nil {
			s.logger.Error(err.Error())
		}
	}()

	return links, nil
}

func (s *link) CreateLink(ctx context.Context, name, addr, icon string, place int) (*types.Link, error) {
	nameTrim := strings.TrimSpace(name)
	addrTrim := strings.TrimSpace(addr)
	if nameTrim == "" || addrTrim == "" {
		return nil, types.NewErrBadRequest(errors.New("link's name or url can't be an empty string"))
	}

	errCh := make(chan error)
	go func() {
		_, err := url.ParseRequestURI(addrTrim)
		errCh <- err
	}()

	link := types.Link{
		Id:      uuid.NewString(),
		Name:    nameTrim,
		URL:     addr,
		Created: time.Now().UTC().Format(time.RFC3339),
		Icon:    icon,
		Place:   place,
	}

	err := <-errCh
	if err != nil {
		return nil, types.NewErrBadRequest(err)
	}

	go func() {
		s.cache.Flush(context.Background())
	}()

	errCh = make(chan error)
	go func(errChan chan error) {
		_, err := s.linkRepo.Create(ctx, link)
		errChan <- err
	}(errCh)

	return &link, <-errCh
}

func (s *link) DeleteLink(ctx context.Context, id string) (*types.Link, error) {
	go func() {
		timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := s.cache.Delete(timeout, id); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return s.linkRepo.Delete(ctx, id)
}

func (s *link) Seed(ctx context.Context) error {
	_, err := s.linkRepo.Create(ctx, types.Link{
		Id:      "09741221-7ea7-4106-ac19-8d2c2c90afbc",
		Name:    "reddit",
		URL:     "https://reddit.com",
		Created: time.Now().Format(time.RFC3339),
		Icon:    "reddit",
		Place:   1,
	})
	if err != nil {
		return err
	}
	_, err = s.linkRepo.Create(ctx, types.Link{
		Id:      "c46428bd-a807-4042-812b-f3b56f047732",
		Name:    "my github",
		URL:     "https://github.com/yosa12978",
		Created: time.Now().Format(time.RFC3339),
		Icon:    "github",
		Place:   0,
	})
	if err != nil {
		return err
	}
	_, err = s.linkRepo.Create(ctx, types.Link{
		Id:      "60a9f6e8-8fda-480a-832a-3e3a07ae8890",
		Name:    "wow forum (icy veins)",
		URL:     "https://www.icy-veins.com/",
		Created: time.Now().Format(time.RFC3339),
		Icon:    "",
		Place:   2,
	})
	return err
}

func (s *link) GetLinkById(ctx context.Context, id string) (*types.Link, error) {
	link, err := s.cache.GetLinkById(ctx, id)
	if err != nil {
		if errors.Is(err, types.ErrInternalFailure) {
			s.logger.Error(err.Error())
		}
	}
	if link != nil && !errors.Is(err, types.ErrNotFound) {
		return link, nil
	}

	link, err = s.linkRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := s.cache.AddLink(timeout, *link); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return link, nil
}
