package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	GetPostComments(ctx context.Context, postId string) (types.Page[types.Comment], int64, error)
	GetCommentById(ctx context.Context, id string) (*types.Comment, error)
	AddPostComments(ctx context.Context, page int, comment ...types.Comment) error
	AddComment(ctx context.Context, comment types.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetCommentsCount(ctx context.Context, postId string) (int, error)
	SetCommentsCount(ctx context.Context, postId string) error
}

type commentRedis struct {
	rdb    *redis.Client
	logger logging.Logger
}

func NewCommentRedis(rdb *redis.Client, logger logging.Logger) Comment {
	return &commentRedis{
		rdb:    rdb,
		logger: logger,
	}
}

func (c *commentRedis) AddComment(ctx context.Context, comment types.Comment) error {
	panic("unimplemented")
}

func (c *commentRedis) AddPostComments(ctx context.Context, page int, comment ...types.Comment) error {
	panic("unimplemented")
}

func (c *commentRedis) DeleteComment(ctx context.Context, id string) error {
	panic("unimplemented")
}

func (c *commentRedis) GetCommentById(ctx context.Context, id string) (*types.Comment, error) {
	panic("unimplemented")
}

func (c *commentRedis) GetCommentsCount(ctx context.Context, postId string) (int, error) {
	panic("unimplemented")
}

func (c *commentRedis) GetPostComments(ctx context.Context, postId string) (types.Page[types.Comment], int64, error) {
	panic("unimplemented")
}

func (c *commentRedis) SetCommentsCount(ctx context.Context, postId string) error {
	panic("unimplemented")
}
