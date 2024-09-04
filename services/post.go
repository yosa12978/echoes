package services

import (
	"context"
	"errors"
	"fmt"
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
	CreatePost(ctx context.Context, title, content string, tweet bool) (*types.Post, error)
	DeletePost(ctx context.Context, id string) (*types.Post, error)
	Seed(ctx context.Context) error
	Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error)
}

type post struct {
	postRepo     repos.Post
	postCache    cache.Post
	logger       logging.Logger
	postSearcher repos.PostSearcher
}

func NewPost(postRepo repos.Post, postCache cache.Post, logger logging.Logger, postSearcher repos.PostSearcher) Post {
	return &post{
		postRepo:     postRepo,
		postCache:    postCache,
		logger:       logger,
		postSearcher: postSearcher,
	}
}

func (s *post) GetPosts(ctx context.Context) ([]types.Post, error) {
	return s.postRepo.FindAll(ctx)
}

func (s *post) GetPostsPaged(ctx context.Context, page, size int) (*types.Page[types.Post], error) {
	pageFromCache, version, err := s.postCache.GetPostsByPage(ctx, page, size)
	if err != nil {
		if errors.Is(err, types.ErrInternalFailure) {
			s.logger.Error(err.Error())
		}
	}
	if pageFromCache != nil {
		return pageFromCache, nil
	}

	t := time.UnixMicro(version).Format(time.RFC3339)
	postsPage, err := s.postRepo.GetPageTime(ctx, t, page, size)
	if err != nil {
		return nil, err
	}

	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.postCache.AddPageOfPosts(timeout, page, *postsPage); err != nil {
			if errors.Is(err, types.ErrInternalFailure) {
				s.logger.Error(err.Error())
			}
		}
	}()

	return postsPage, err
}

func (s *post) GetPostById(ctx context.Context, id string) (*types.Post, error) {
	postFromCache, err := s.postCache.GetPostById(ctx, id)
	if err == nil {
		return postFromCache, nil
	}
	if errors.Is(err, types.ErrInternalFailure) {
		s.logger.Error(err.Error())
	}

	post, err := s.postRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.postCache.AddPost(ctx, *post); err != nil {
			s.logger.Error(err.Error())
		}
	}()

	return post, nil
}

func (s *post) PinPost(ctx context.Context, id string) (*types.Post, error) {
	post, err := s.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}
	post.Pinned = !post.Pinned

	go func() {
		if err := s.postCache.PinPost(ctx, id); err != nil {
			if errors.Is(err, types.ErrInternalFailure) {
				s.logger.Error(err.Error())
			}
		}
	}()

	return s.postRepo.Update(ctx, post.Id, *post)
}

func (s *post) CreatePost(ctx context.Context, title, content string, tweet bool) (*types.Post, error) {
	titleTrim := strings.TrimSpace(title)
	contentTrim := strings.TrimSpace(content)
	if titleTrim == "" || contentTrim == "" {
		return nil, types.NewErrBadRequest(errors.New("can't create post with empty title or content"))
	}
	id := uuid.NewString()
	post := types.Post{
		Id:      id,
		Title:   titleTrim,
		Content: contentTrim,
		Created: time.Now().UTC().Format(time.RFC3339),
		Pinned:  false,
		Tweet:   tweet,
	}

	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.postCache.AddPost(timeout, post); err != nil {
			s.logger.Error(err.Error())
		}
	}()

	return s.postRepo.Create(ctx, post)
}

func (s *post) DeletePost(ctx context.Context, id string) (*types.Post, error) {
	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.postCache.Delete(timeout, id); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return s.postRepo.Delete(ctx, id)
}

func (s *post) Seed(ctx context.Context) error {
	for i := 0; i < 60; i++ {
		time.Sleep(1000 * time.Millisecond)
		_, err := s.postRepo.Create(ctx, types.NewPost(fmt.Sprintf("post #%d", 30-i), fmt.Sprintf("post content #%d", 30-i), false))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *post) Search(ctx context.Context, query string, page, size int) (*types.Page[types.Post], error) {
	return s.postSearcher.Search(ctx, query, page, size)
}
