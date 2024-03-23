package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	GetPostComments(ctx context.Context, postId string, page, size int) (*types.CommentsInfo, error)
	GetCommentById(ctx context.Context, commentId string) (*types.Comment, error)
	CreateComment(ctx context.Context, postId, name, email, content string) (*types.Comment, error)
	DeleteComment(ctx context.Context, commentId string) (*types.Comment, error)
	GetCommentsCount(ctx context.Context, postId string) (int, error)
	Seed(ctx context.Context) error
}

type comment struct {
	commentRepo repos.Comment
	postService Post
	cache       cache.Cache
	logger      logging.Logger
}

func NewComment(commentRepo repos.Comment, postService Post, cache cache.Cache, logger logging.Logger) Comment {
	return &comment{
		commentRepo: commentRepo,
		postService: postService,
		cache:       cache,
		logger:      logger,
	}
}

func (s *comment) GetPostComments(ctx context.Context, postId string, page, size int) (*types.CommentsInfo, error) {
	if _, err := s.postService.GetPostById(ctx, postId); err != nil {
		return nil, err
	}
	commentsPaged, err := s.commentRepo.GetPage(ctx, postId, page, size)
	if err != nil {
		// log err here
		return nil, err
	}
	res := types.CommentsInfo{
		Page:   *commentsPaged,
		PostId: postId,
	}
	return &res, nil
}

func (s *comment) GetCommentById(ctx context.Context, commentId string) (*types.Comment, error) {
	commentFromCache, err := s.cache.Get(ctx, "comments:"+commentId)
	if err == nil {
		var comment types.Comment
		err := json.Unmarshal([]byte(commentFromCache), &comment)
		return &comment, err
	}

	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return nil, err
	}

	go func() {
		commentJson, _ := json.Marshal(comment)
		_, err := s.cache.Set(ctx, "comments:"+commentId, commentJson, 0)
		if err != nil {
			s.logger.Error(err)
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
		Created: time.Now().Format(time.RFC3339),
		Name:    name,
		Email:   email,
		Content: content,
		PostId:  postId,
	}

	go func() {
		commentJson, _ := json.Marshal(comm)
		_, err := s.cache.Set(ctx, "comments:"+comm.Id, commentJson, 0)
		if err != nil {
			s.logger.Error(err)
		}
	}()
	return s.commentRepo.Create(ctx, comm)
}

func (s *comment) DeleteComment(ctx context.Context, commentId string) (*types.Comment, error) {
	go func() {
		_, err := s.cache.Del(ctx, "comments:"+commentId)
		if err != nil {
			s.logger.Error(err)
		}
	}()
	return s.commentRepo.Delete(ctx, commentId)
}

func (s *comment) Seed(ctx context.Context) error {
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("Name#%d", i)
		email := fmt.Sprintf("email%d@email.com", i)
		content := fmt.Sprintf("content %d", time.Now().UnixNano())
		_, err := s.CreateComment(ctx, "895cef0a-58e0-4f55-b49f-6bea42d8bcd1", name, email, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *comment) GetCommentsCount(ctx context.Context, postId string) (int, error) {
	countstr, err := s.cache.Get(ctx, "comments_count:"+postId)
	if err == nil {
		return strconv.Atoi(countstr)
	}

	count, err := s.commentRepo.GetCommentsCount(ctx, postId)
	if err != nil {
		return 0, err
	}
	go func() {
		_, err = s.cache.Set(ctx, "comments_count:"+postId, count, 0)
		if err != nil {
			s.logger.Error(err)
		}
	}()
	return count, nil
}
