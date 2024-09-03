package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	GetPostComments(ctx context.Context, postId string, page, size int) (*types.Page[types.Comment], error)
	GetCommentById(ctx context.Context, commentId string) (*types.Comment, error)
	CreateComment(ctx context.Context, postId, name, email, content string) (*types.Comment, error)
	DeleteComment(ctx context.Context, commentId string) (*types.Comment, error)
	GetCommentsCount(ctx context.Context, postId string) (int, error)
	Seed(ctx context.Context) error
}

type comment struct {
	commentRepo repos.Comment
	postService Post
	cache       cache.Comment
	logger      logging.Logger
}

func NewComment(commentRepo repos.Comment, postService Post, cache cache.Comment, logger logging.Logger) Comment {
	return &comment{
		commentRepo: commentRepo,
		postService: postService,
		cache:       cache,
		logger:      logger,
	}
}

func (s *comment) GetPostComments(ctx context.Context, postId string, page, size int) (*types.Page[types.Comment], error) {
	// if _, err := s.postService.GetPostById(ctx, postId); err != nil {
	// 	return nil, err
	// }
	commentsFromCache, version, err := s.cache.GetPostComments(ctx, postId, page)
	if err != nil {
		if errors.Is(err, cache.ErrInternalFailure) {
			s.logger.Error(err.Error())
		}
	}
	if commentsFromCache != nil { // cheching err again because ErrNotFound might occur. Refactor this later
		return commentsFromCache, nil
	}

	t := time.UnixMicro(version).Format(time.RFC3339)
	commentsPaged, err := s.commentRepo.GetPageTime(ctx, t, postId, page, size)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.cache.AddPostComments(timeout, postId, page, *commentsPaged)
	}()
	return commentsPaged, nil
}

func (s *comment) GetCommentById(ctx context.Context, commentId string) (*types.Comment, error) {
	commentFromCache, err := s.cache.GetCommentById(ctx, commentId)
	if err == nil {
		return commentFromCache, nil
	}
	if errors.Is(err, cache.ErrInternalFailure) {
		s.logger.Error(err.Error())
	}

	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return nil, err
	}

	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.cache.AddComment(timeout, *comment); err != nil {
			s.logger.Error(err.Error())
		}
	}()

	return comment, nil
}

func (s *comment) CreateComment(ctx context.Context, postId, name, email, content string) (*types.Comment, error) {
	if _, err := s.postService.GetPostById(ctx, postId); err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	content = strings.TrimSpace(content)
	if name == "" || email == "" || content == "" {
		return nil, errors.New("name, email or content field can't be empty")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("email address isn't valid")
	}
	comm := types.Comment{
		Id:      uuid.NewString(),
		Created: time.Now().UTC().Format(time.RFC3339),
		Name:    name,
		Email:   email,
		Content: content,
		PostId:  postId,
	}

	go func() { // i don't like this. refactor
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.cache.AddComment(timeout, comm); err != nil {
			s.logger.Error(err.Error())
		}
		if _, err := s.cache.RefreshPagination(timeout, postId); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return s.commentRepo.Create(ctx, comm)
}

func (s *comment) DeleteComment(ctx context.Context, commentId string) (*types.Comment, error) {
	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.cache.DeleteComment(timeout, commentId); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return s.commentRepo.Delete(ctx, commentId)
}

func (s *comment) Seed(ctx context.Context) error {
	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		name := fmt.Sprintf("Name#%d", i)
		email := fmt.Sprintf("email%d@email.com", i)
		content := fmt.Sprintf("content %d", time.Now().UnixNano())
		_, err := s.CreateComment(ctx, "dcc0650c-4370-4ef3-a846-a7d71ddd55fc", name, email, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *comment) GetCommentsCount(ctx context.Context, postId string) (int, error) {
	count, err := s.cache.GetCommentsCount(ctx, postId)
	if err == nil {
		return count, nil
	}

	count, err = s.commentRepo.GetCommentsCount(ctx, postId)
	if err != nil {
		return 0, err
	}
	go func() {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.cache.SetCommentsCount(timeout, postId, count); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return count, nil
}
