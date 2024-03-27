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

func (s *comment) getPaginationVersion(ctx context.Context, postId string) (int64, error) {
	var version int64
	key := "comments_pagination_version:" + postId
	versionFromCache, err := s.cache.Get(ctx, key)
	if err != nil {
		version := time.Now().UnixMicro()
		_, err = s.cache.Set(
			ctx,
			key,
			version,
			1*time.Minute,
		)
		return version, err
	}
	version, err = strconv.ParseInt(versionFromCache, 10, 64)
	return version, err
}

func (s *comment) GetPostComments(ctx context.Context, postId string, page, size int) (*types.CommentsInfo, error) {
	if _, err := s.postService.GetPostById(ctx, postId); err != nil {
		return nil, err
	}
	version, err := s.getPaginationVersion(ctx, postId)
	if err != nil {
		s.logger.Error(err)
	}
	key := fmt.Sprintf("comments:%s:%v:%d", postId, version, page)
	commentsFromCache, err := s.cache.Get(ctx, key)
	if err == nil {
		var res types.CommentsInfo
		err := json.Unmarshal([]byte(commentsFromCache), &res)
		if err == nil {
			return &res, nil
		}
		s.logger.Error(err)
	}

	t := time.UnixMicro(version).Format(time.RFC3339)
	commentsPaged, err := s.commentRepo.GetPageTime(ctx, t, postId, page, size)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	res := types.CommentsInfo{
		Page:   *commentsPaged,
		PostId: postId,
	}

	go func() {
		pageJson, _ := json.Marshal(res)
		s.cache.SetNX(ctx, key, pageJson, 65*time.Second)
	}()
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
	countstr, err := s.cache.Get(ctx, "comments_count:"+postId)
	if err == nil {
		return strconv.Atoi(countstr)
	}

	count, err := s.commentRepo.GetCommentsCount(ctx, postId)
	if err != nil {
		return 0, err
	}
	go func() {
		_, err = s.cache.Set(ctx, "comments_count:"+postId, count, 60*time.Second)
		if err != nil {
			s.logger.Error(err)
		}
	}()
	return count, nil
}
