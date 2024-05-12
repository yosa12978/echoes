package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Link interface {
	GetLinks(ctx context.Context) ([]types.Link, error)
	CreateLink(ctx context.Context, name, url, icon string) (*types.Link, error)
	DeleteLink(ctx context.Context, id string) (*types.Link, error)
	Seed(ctx context.Context) error
}

type link struct {
	linkRepo repos.Link
	cache    cache.Cache
	logger   logging.Logger
}

func NewLink(linkRepo repos.Link, cache cache.Cache, logger logging.Logger) Link {
	return &link{linkRepo: linkRepo, cache: cache, logger: logger}
}

// Migrate to hashmaps instead of json-encoded string for caching
// Here is some helpful funcs
// func linkToMap(l types.Link) map[string]interface{} {
// 	return map[string]interface{}{
// 		"Id":      l.Id,
// 		"Name":    l.Name,
// 		"Created": l.Created,
// 		"URL":     l.URL,
// 	}
// }

// func mapToLink(m map[string]string) types.Link {
// 	return types.Link{
// 		Id:      m["Id"],
// 		Name:    m["Name"],
// 		Created: m["Created"],
// 		URL:     m["URL"],
// 	}
// }

func (s *link) GetLinks(ctx context.Context) ([]types.Link, error) {
	inCache, _ := s.cache.Exists(ctx, "links")
	if inCache == 1 {
		setRes, err := s.cache.ZRange(ctx, "links", 0, -1)
		if err == nil {
			links := make([]types.Link, len(setRes))
			linksFromSet, _ := s.cache.MGet(ctx, setRes...)
			var wg sync.WaitGroup
			for k, v := range linksFromSet {
				wg.Add(1)
				go func(key int, val string) {
					defer wg.Done()
					json.Unmarshal([]byte(val), &links[key])
				}(k, v.(string))
			}
			wg.Wait()
			return links, nil
		}
	}

	links, err := s.linkRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		tx, _ := s.cache.Tx()
		tx.Append(ctx, func(pipe cache.Cache) error {
			members := make([]cache.Member, len(links))
			for k, v := range links {
				link, _ := json.Marshal(v)
				key := fmt.Sprintf("links:%s", v.Id)
				pipe.Set(ctx, key, link, 0)
				members[k] = cache.Member{Member: key, Score: float64(k)}
			}
			pipe.ZAdd(ctx, "links", members...)
			return nil
		})
		tx.Exec(ctx)
	}()

	return links, err
}

func (s *link) CreateLink(ctx context.Context, name, addr, icon string) (*types.Link, error) {
	nameTrim := strings.TrimSpace(name)
	addrTrim := strings.TrimSpace(addr)
	if nameTrim == "" || addrTrim == "" {
		return nil, errors.New("link's name or url can't be an empty string")
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
		Created: time.Now().Format(time.RFC3339),
		Icon:    icon,
	}

	err := <-errCh
	if err != nil {
		return nil, err
	}

	go func() {
		linkJson, _ := json.Marshal(link)
		linkKey := fmt.Sprintf("links:%s", link.Id)
		tx, _ := s.cache.Tx()
		tx.Append(ctx, func(pipe cache.Cache) error {
			pipe.Set(ctx, linkKey, string(linkJson), 0)
			score, _ := s.cache.ZCard(ctx, "posts")
			pipe.ZAdd(ctx, "links", cache.Member{Score: float64(score + 1), Member: linkKey})
			return nil
		})
		tx.Exec(ctx)
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
		key := fmt.Sprintf("links:%s", id)
		tx, _ := s.cache.Tx()
		tx.Append(ctx, func(pipe cache.Cache) error {
			_, err := pipe.Del(ctx, key)
			pipe.ZRem(ctx, "links", key)
			return err
		})
		tx.Exec(ctx)
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
	})
	return err
}
