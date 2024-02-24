package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Post interface {
	GetPosts(ctx context.Context) []types.Post
	GetPostsPaged(ctx context.Context, page, size int) *types.Page[types.Post]
	GetPostById(ctx context.Context, id string) (*types.Post, error)
	// pin post works like a trigger
	PinPost(ctx context.Context, id string) (*types.Post, error)
	CreatePost(ctx context.Context, title, content string) (*types.Post, error)
	DeletePost(ctx context.Context, id string) (*types.Post, error)
	Seed(ctx context.Context) error
}

type post struct {
	postRepo repos.Post
}

func NewPost(postRepo repos.Post) Post {
	return &post{postRepo: postRepo}
}

func (s *post) GetPosts(ctx context.Context) []types.Post {
	return s.postRepo.FindAll(ctx)
}

func (s *post) GetPostsPaged(ctx context.Context, page, size int) *types.Page[types.Post] {
	return s.postRepo.GetPage(ctx, page, size)
}

func (s *post) GetPostById(ctx context.Context, id string) (*types.Post, error) {
	return s.postRepo.FindById(ctx, id)
}

func (s *post) PinPost(ctx context.Context, id string) (*types.Post, error) {
	post, err := s.postRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	post.Pinned = !post.Pinned
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
	return s.postRepo.Create(ctx, post)
}

func (s *post) DeletePost(ctx context.Context, id string) (*types.Post, error) {
	return s.postRepo.Delete(ctx, id)
}

func (s *post) Seed(ctx context.Context) error {
	for i := 0; i < 50; i++ {
		time.Sleep(50 * time.Millisecond)
		_, err := s.postRepo.Create(ctx, types.NewPost(fmt.Sprintf("post #%d", i), fmt.Sprintf("post content #%d", i)))
		if err != nil {
			return err
		}
	}
	return nil
}
