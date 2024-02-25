package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	GetPostComments(ctx context.Context, postId string, page, size int) (*types.CommentsInfo, error)
	GetCommentById(ctx context.Context, commentId string) (*types.Comment, error)
	CreateComment(ctx context.Context, postId, name, email, content string) (*types.Comment, error)
	DeleteComment(ctx context.Context, commentId string) (*types.Comment, error)
}

type comment struct {
	commentRepo repos.Comment
	postService Post
}

func NewComment(commentRepo repos.Comment, postService Post) Comment {
	return &comment{commentRepo: commentRepo, postService: postService}
}

func (s *comment) GetPostComments(ctx context.Context, postId string, page, size int) (*types.CommentsInfo, error) {
	if _, err := s.postService.GetPostById(ctx, postId); err != nil {
		return nil, err
	}
	commentsPaged := s.commentRepo.GetPage(ctx, postId, page, size)
	res := types.CommentsInfo{
		Page:   *commentsPaged,
		PostId: postId,
	}
	return &res, nil
}

func (s *comment) GetCommentById(ctx context.Context, commentId string) (*types.Comment, error) {
	return s.commentRepo.FindById(ctx, commentId)
}

func (s *comment) CreateComment(ctx context.Context, postId, name, email, content string) (*types.Comment, error) {
	if _, err := s.commentRepo.FindById(ctx, postId); err != nil {
		return nil, errors.New("post not found")
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
	return s.commentRepo.Create(ctx, comm)
}

func (s *comment) DeleteComment(ctx context.Context, commentId string) (*types.Comment, error) {
	return s.commentRepo.Delete(ctx, commentId)
}
