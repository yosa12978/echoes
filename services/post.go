package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Post interface {
	GetPosts(ctx context.Context) ([]types.Post, error)
	GetPostsPaged(ctx context.Context, page, size int) (*types.Page[types.Post], error)
	GetPostById(ctx context.Context, id string) (*types.Post, error)
	// pin post works like a trigger
	PinPost(ctx context.Context, id string) (*types.Post, error)
	CreatePost(ctx context.Context, title, content string) (*types.Post, error)
	DeletePost(ctx context.Context, id string) (*types.Post, error)
	Seed(ctx context.Context) error
	Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error)
}

type post struct {
	postRepo repos.Post
	cache    cache.Cache
	logger   logging.Logger
}

func NewPost(postRepo repos.Post, cache cache.Cache, logger logging.Logger) Post {
	return &post{
		postRepo: postRepo,
		cache:    cache,
		logger:   logger,
	}
}

func (s *post) GetPosts(ctx context.Context) ([]types.Post, error) {
	return s.postRepo.FindAll(ctx)
}

func (s *post) getPaginationVersion(ctx context.Context) (int64, error) {
	var version int64
	versionFromCache, err := s.cache.Get(ctx, "posts_pagination_version")
	if err != nil {
		version := time.Now().UnixMicro()
		_, err = s.cache.Set(
			ctx,
			"posts_pagination_version",
			version,
			1*time.Minute,
		)
		return version, err
	}
	version, err = strconv.ParseInt(versionFromCache, 10, 64)
	return version, err
}

func (s *post) GetPostsPaged(ctx context.Context, page, size int) (*types.Page[types.Post], error) {
	version, _ := s.getPaginationVersion(ctx)
	key := fmt.Sprintf("posts:%v:page:%d", version, page)
	pageFromCache, err := s.cache.Get(ctx, key)
	if err == nil {
		var res *types.Page[types.Post]
		err := json.Unmarshal([]byte(pageFromCache), &res)
		if err == nil {
			return res, err
		}
		s.logger.Error(err)
	}

	t := time.UnixMicro(version).Format(time.RFC3339)
	postsPage, err := s.postRepo.GetPageTime(ctx, t, page, size)
	if err != nil {
		return nil, err
	}

	go func() {
		pageJson, _ := json.Marshal(postsPage)
		_, err := s.cache.SetNX(ctx, key, pageJson, 65*time.Second)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return postsPage, err
}

func (s *post) GetPostById(ctx context.Context, id string) (*types.Post, error) {
	postJson, err := s.cache.Get(ctx, "posts:"+id)
	if err == nil {
		var post types.Post
		json.Unmarshal([]byte(postJson), &post)
		return &post, nil
	}

	post, err := s.postRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		postBytes, _ := json.Marshal(post)
		_, err = s.cache.Set(ctx, "posts:"+id, string(postBytes), 0)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return post, err
}

func (s *post) PinPost(ctx context.Context, id string) (*types.Post, error) {
	post, err := s.postRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	post.Pinned = !post.Pinned

	go func() {
		postb, _ := json.Marshal(post)
		_, err := s.cache.SetXX(ctx, "posts:"+id, string(postb), 0)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return s.postRepo.Update(ctx, post.Id, *post)
}

func (s *post) CreatePost(ctx context.Context, title, content string) (*types.Post, error) {
	titleTrim := strings.TrimSpace(title)
	contentTrim := strings.TrimSpace(content)
	if titleTrim == "" || contentTrim == "" {
		return nil, errors.New("can't create post with empty title or content")
	}
	post := types.Post{
		Id:      uuid.NewString(),
		Title:   titleTrim,
		Content: contentTrim,
		Created: time.Now().Format(time.RFC3339),
		Pinned:  false,
	}

	go func() {
		postb, _ := json.Marshal(post)
		key := "posts:" + post.Id
		tx, _ := s.cache.Tx()
		tx.Append(ctx, func(pipe cache.Cache) error {
			s.cache.Set(ctx, key, string(postb), 0)
			score, _ := s.cache.ZCard(ctx, "posts")
			s.cache.ZAdd(ctx, "posts", cache.Member{
				Score:  float64(score + 1),
				Member: key,
			})
			return nil
		})
		if err := tx.Exec(ctx); err != nil {
			s.logger.Error(err)
		}
	}()

	return s.postRepo.Create(ctx, post)
}

func (s *post) DeletePost(ctx context.Context, id string) (*types.Post, error) {
	go func() {
		tx, _ := s.cache.Tx()
		key := "posts:" + id
		tx.Append(ctx, func(pipe cache.Cache) error {
			_, err := pipe.Del(ctx, key)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			_, err = pipe.ZRem(ctx, "posts", key)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			return nil
		})
		if err := tx.Exec(ctx); err != nil {
			s.logger.Error(err)
		}
	}()
	return s.postRepo.Delete(ctx, id)
}

func (s *post) Seed(ctx context.Context) error {
	for i := 0; i < 30; i++ {
		time.Sleep(1000 * time.Millisecond)
		_, err := s.postRepo.Create(ctx, types.NewPost(fmt.Sprintf("post #%d", 30-i), fmt.Sprintf("post content #%d", 30-i)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *post) Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error) {
	return nil, nil
}
