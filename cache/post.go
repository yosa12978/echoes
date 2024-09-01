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

type Post interface {
	GetPostsByPage(ctx context.Context, page, size int) (*types.Page[types.Post], int64, error)
	GetPostById(ctx context.Context, id string) (*types.Post, error)

	PinPost(ctx context.Context, id string) error
	AddPost(ctx context.Context, post types.Post) error
	AddPageOfPosts(ctx context.Context, pageNum int, page types.Page[types.Post]) error
	Update(ctx context.Context, id string, post types.Post) error
	Delete(ctx context.Context, id string) error
}

type postRedis struct {
	rdb    *redis.Client
	logger logging.Logger
}

func NewPostRedis(rdb *redis.Client, logger logging.Logger) Post {
	return &postRedis{rdb: rdb, logger: logger}
}

func (p *postRedis) GetPostById(ctx context.Context, id string) (*types.Post, error) {
	postMap, err := p.rdb.HGetAll(ctx, "posts:"+id).Result()
	if len(postMap) == 0 {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	pinned, _ := strconv.ParseBool(postMap["pinned"])
	tweet, _ := strconv.ParseBool(postMap["tweet"])
	comments, _ := strconv.Atoi(postMap["comments"])
	post := types.Post{
		Id:       postMap["id"],
		Title:    postMap["title"],
		Content:  postMap["content"],
		Created:  postMap["created"],
		Pinned:   pinned,
		Tweet:    tweet,
		Comments: comments,
	}
	return &post, nil
}

func (p *postRedis) refreshPaginationVersion(ctx context.Context) (int64, error) {
	version := time.Now().UnixMicro()
	res, err := p.rdb.Set(ctx, "posts_pagination_version", version, 1*time.Minute).Result()
	if res != "OK" || err != nil {
		return 0, fmt.Errorf("failed to update posts_pagination_version: %w", ErrInternalFailure)
	}
	p.logger.Info("updated posts_pagination_version", "version", version)
	return version, err
}

func (p *postRedis) getPaginationVersion(ctx context.Context) (int64, error) {
	versionFromCache, err := p.rdb.Get(ctx, "posts_pagination_version").Result()
	if err != nil {
		return p.refreshPaginationVersion(ctx)
	}
	return strconv.ParseInt(versionFromCache, 10, 64)
}

// latency here 2ms+ size=20
func (p *postRedis) GetPostsByPage(ctx context.Context, page int, size int) (*types.Page[types.Post], int64, error) {
	version, _ := p.getPaginationVersion(ctx)
	valueKey := fmt.Sprintf("posts:%v:page:%d", version, page)
	metaKey := fmt.Sprintf("posts:%v:page_meta:%d", version, page)
	//make this concurrent
	valueExists, _ := p.rdb.Exists(ctx, valueKey).Result()
	metaExists, _ := p.rdb.Exists(ctx, metaKey).Result()
	if valueExists == 0 || metaExists == 0 {
		return nil, version, ErrNotFound
	}
	postKeysJson, err := p.rdb.Get(ctx, valueKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, version, nil
		}
		return nil, version, errors.Join(err, ErrInternalFailure)
	}
	posts := make([]types.Post, 0, len(postKeysJson))
	var postIDs []string
	json.Unmarshal([]byte(postKeysJson), &postIDs)
	for _, postId := range postIDs {
		// make this concurrent
		post, err := p.GetPostById(ctx, postId)
		if err != nil {
			return nil, version, err
		}
		posts = append(posts, *post)
	}
	// check for key existance here instead of beginning
	// or even move this part to the beginning before posts fetching
	metaMap, err := p.rdb.HGetAll(ctx, metaKey).Result()
	if len(metaMap) == 0 {
		return nil, version, ErrNotFound
	}
	if err != nil {
		return nil, version, errors.Join(err, ErrInternalFailure)
	}
	hasNext, _ := strconv.ParseBool(metaMap["has_next"])
	page_size, _ := strconv.Atoi(metaMap["size"])
	total, _ := strconv.Atoi(metaMap["total"])
	nextPage, _ := strconv.Atoi(metaMap["next_page"])
	postsPage := types.Page[types.Post]{
		Content:  posts,
		HasNext:  hasNext,
		Size:     page_size,
		Total:    total,
		NextPage: nextPage,
	}
	return &postsPage, version, nil
}

// combine this with addPostPipeline
func (p *postRedis) AddPost(ctx context.Context, post types.Post) error {
	pipe := p.rdb.Pipeline()

	if err := addPost(ctx, post, pipe); err != nil {
		return err
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	_, err := p.refreshPaginationVersion(ctx) // replace this with adding to cache instead
	return err
}

func addPost(ctx context.Context, post types.Post, rdb redis.Cmdable) error {
	postMap := map[string]interface{}{
		"id":       post.Id,
		"title":    post.Title,
		"content":  post.Content,
		"created":  post.Created,
		"pinned":   post.Pinned,
		"tweet":    post.Tweet,
		"comments": post.Comments,
	}
	key := fmt.Sprintf("posts:%s", post.Id)
	_, err := rdb.HSet(ctx, key, postMap).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	_, err = rdb.Expire(ctx, key, 150*time.Second).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	return err
}

func (p *postRedis) AddPageOfPosts(
	ctx context.Context,
	pageNum int,
	page types.Page[types.Post],
) error {
	version, _ := p.getPaginationVersion(ctx)

	pipe := p.rdb.Pipeline()

	// caching posts and preparing data for sorted set
	// make this concurrent
	postsIDs := make([]string, 0, len(page.Content))
	for _, val := range page.Content {
		if err := addPost(ctx, val, pipe); err != nil {
			return err
		}
		postsIDs = append(postsIDs, val.Id)
	}
	postsKeysJson, _ := json.Marshal(postsIDs)

	// caching posts sorted set
	valueKey := fmt.Sprintf("posts:%v:page:%d", version, pageNum)
	err := pipe.Set(ctx, valueKey, postsKeysJson, 2*time.Minute).Err()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}

	// setting up metadata
	metaKey := fmt.Sprintf("posts:%v:page_meta:%d", version, pageNum)
	metaData := map[string]interface{}{
		"has_next":  page.HasNext,
		"next_page": page.NextPage,
		"total":     page.Total,
		"size":      page.Size,
	}
	_, err = pipe.HSet(ctx, metaKey, metaData).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	err = pipe.Expire(ctx, metaKey, 2*time.Minute).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}

	_, err = pipe.Exec(ctx)

	return err
}

func (p *postRedis) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf("posts:%s", id)
	_, err := p.rdb.Del(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	_, err = p.refreshPaginationVersion(ctx)
	return err
}

func (p *postRedis) PinPost(ctx context.Context, id string) error {
	key := fmt.Sprintf("posts:%s", id)

	pinnedStr, err := p.rdb.HGet(ctx, key, "pinned").Result()
	if len(pinnedStr) == 0 {
		return ErrNotFound
	}
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	pinned, _ := strconv.ParseBool(pinnedStr)
	pinned = !pinned
	_, err = p.rdb.HSet(ctx, key, "pinned", pinned).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	_, err = p.refreshPaginationVersion(ctx)
	return err
}

func (p *postRedis) Update(ctx context.Context, id string, post types.Post) error {
	panic("unimplemented")
}
