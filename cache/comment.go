package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/types"
)

type Comment interface {
	GetPostComments(ctx context.Context, postId string, page int) (*types.Page[types.Comment], int64, error)
	GetCommentById(ctx context.Context, id string) (*types.Comment, error)
	AddPostComments(ctx context.Context, postId string, page int, comment types.Page[types.Comment]) error
	AddComment(ctx context.Context, comment types.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetCommentsCount(ctx context.Context, postId string) (int, error)
	SetCommentsCount(ctx context.Context, postId string, count int) error
	RefreshPagination(ctx context.Context, postId string) (int64, error)
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

func (c *commentRedis) getPaginationVersion(ctx context.Context, postId string) (int64, error) {
	key := "comments_pagination_version:" + postId
	versionFromCache, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return c.refreshPaginationVersion(ctx, postId)
	}
	return strconv.ParseInt(versionFromCache, 10, 64)
}

func (c *commentRedis) refreshPaginationVersion(ctx context.Context, postId string) (int64, error) {
	key := fmt.Sprintf("comments_pagination_version:%s", postId)
	version := time.Now().UnixMicro()
	res, err := c.rdb.Set(ctx, key, version, 1*time.Minute).Result()
	if res != "OK" || err != nil {
		return 0, fmt.Errorf("failed to update posts_pagination_version: %w", ErrInternalFailure)
	}
	c.logger.Info("updated posts_pagination_version", "version", version)
	return version, err
}

func (c *commentRedis) RefreshPagination(ctx context.Context, postId string) (int64, error) {
	return c.refreshPaginationVersion(ctx, postId)
}

func (c *commentRedis) AddComment(ctx context.Context, comment types.Comment) error {
	key := fmt.Sprintf("comments:%s", comment.Id)
	commentJson, _ := json.Marshal(comment)
	err := c.rdb.Set(ctx, key, commentJson, 1*time.Minute).Err()
	if err != nil {
		return newInternalFailure(err)
	}
	return nil
}

func (c *commentRedis) AddPostComments(
	ctx context.Context,
	postId string,
	page int,
	comments types.Page[types.Comment],
) error {
	version, err := c.getPaginationVersion(ctx, postId)
	if err != nil {
		c.logger.Error(err.Error())
	}
	key := fmt.Sprintf("comments:%s:%v:%d", postId, version, page)
	pageJson, _ := json.Marshal(comments)
	err = c.rdb.SetNX(ctx, key, pageJson, 1*time.Minute).Err()
	if err != nil {
		return newInternalFailure(err)
	}
	return nil
}

func (c *commentRedis) GetPostComments(ctx context.Context, postId string, page int) (*types.Page[types.Comment], int64, error) {
	version, err := c.getPaginationVersion(ctx, postId)
	if err != nil {
		return nil, 0, err
	}
	key := fmt.Sprintf("comments:%s:%v:%d", postId, version, page)
	pageJson, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, 0, ErrNotFound
		}
		return nil, 0, newInternalFailure(err)
	}
	var res types.Page[types.Comment]
	json.Unmarshal([]byte(pageJson), &res)
	return &res, version, nil
}

func (c *commentRedis) DeleteComment(ctx context.Context, id string) error {
	key := fmt.Sprintf("comments:%s", id)
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return newInternalFailure(err)
	}
	return nil
}

func (c *commentRedis) GetCommentById(ctx context.Context, id string) (*types.Comment, error) {
	key := fmt.Sprintf("comments:%s", id)
	commentJson, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, newInternalFailure(err)
	}
	var comment types.Comment
	json.Unmarshal([]byte(commentJson), &comment)
	return &comment, nil
}

// get rid of this
func (c *commentRedis) GetCommentsCount(ctx context.Context, postId string) (int, error) {
	countStr, err := c.rdb.Get(ctx, "comments_count"+postId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrNotFound
		}
		return 0, newInternalFailure(err)
	}
	return strconv.Atoi(countStr)
}

// get rid of this
func (c *commentRedis) SetCommentsCount(ctx context.Context, postId string, count int) error {
	err := c.rdb.Set(ctx, "comments_count:"+postId, count, 60*time.Second).Err()
	if err != nil {
		return newInternalFailure(err)
	}
	return nil
}
