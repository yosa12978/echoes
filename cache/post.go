package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yosa12978/echoes/types"
)

type Post interface {
	GetPostsByPage(ctx context.Context, page, size int) (*types.Page[types.Post], error)
	GetPostById(ctx context.Context, id string) (*types.Post, error)

	PinPost(ctx context.Context, id string) error
	AddPost(ctx context.Context, post types.Post) error
	AddPageOfPosts(ctx context.Context, pageNum int, page types.Page[types.Post]) error
	Update(ctx context.Context, id string, post types.Post) error
	Delete(ctx context.Context, id string) error
}

type postRedis struct {
	rdb *redis.Client
}

func NewPostRedis(rdb *redis.Client) Post {
	return &postRedis{rdb: rdb}
}

func (p *postRedis) getPostByRedisKey(ctx context.Context, key string) (*types.Post, error) {
	postMap, err := p.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
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

func (p *postRedis) GetPostById(ctx context.Context, id string) (*types.Post, error) {
	return p.getPostByRedisKey(ctx, "posts:"+id)
}

func (p *postRedis) refreshPaginationVersion(ctx context.Context) (int64, error) {
	version := time.Now().UnixMicro()
	_, err := p.rdb.Set(ctx, "posts_pagination_version", version, 1*time.Minute).Result()
	return version, err
}

func (p *postRedis) getPaginationVersion(ctx context.Context) (int64, error) {
	versionFromCache, err := p.rdb.Get(ctx, "posts_pagination_version").Result()
	if err != nil {
		return p.refreshPaginationVersion(ctx)
	}
	return strconv.ParseInt(versionFromCache, 10, 64)
}

func (p *postRedis) GetPostsByPage(ctx context.Context, page int, size int) (*types.Page[types.Post], error) {
	version, _ := p.getPaginationVersion(ctx)
	valueKey := fmt.Sprintf("posts:%v:page:%d", version, page)
	metaKey := fmt.Sprintf("posts:%v:page_meta:%d", version, page)
	//make this concurrent
	valueExists, _ := p.rdb.Exists(ctx, valueKey).Result()
	metaExists, _ := p.rdb.Exists(ctx, metaKey).Result()
	if valueExists == 0 || metaExists == 0 {
		return nil, ErrNotFound
	}
	postKeys, err := p.rdb.ZRange(ctx, valueKey, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	posts := make([]types.Post, 0, len(postKeys))
	for _, postKey := range postKeys {
		// make this concurrent
		post, err := p.getPostByRedisKey(ctx, postKey)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}
	// check for key existance here instead of beginning
	// or even move this part to the beginning before posts fetching
	metaMap, err := p.rdb.HGetAll(ctx, metaKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
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
	return &postsPage, nil
}

func (p *postRedis) AddPost(ctx context.Context, post types.Post) error {
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
	_, err := p.rdb.HSet(ctx, key, postMap).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	_, err = p.rdb.Expire(ctx, key, 90*time.Second).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	return nil
}

func (p *postRedis) AddPageOfPosts(ctx context.Context, pageNum int, page types.Page[types.Post]) error {
	version, _ := p.getPaginationVersion(ctx)
	metaKey := fmt.Sprintf("posts:%v:page_meta:%d", version, pageNum)
	metaData := map[string]interface{}{
		"has_next":  page.HasNext,
		"next_page": page.HasNext,
		"total":     page.Total,
		"size":      page.Size,
	}
	_, err := p.rdb.HSet(ctx, metaKey, metaData).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	postsKeys := make([]redis.Z, 0, len(page.Content))
	for _, val := range page.Content {
		if err := p.AddPost(ctx, val); err != nil {
			return err
		}
		timestamp, _ := time.Parse(time.RFC3339, val.Created)
		postsKeys = append(postsKeys,
			redis.Z{
				Member: "posts:" + val.Id,
				Score:  float64(timestamp.Unix()),
			})
	}
	valueKey := fmt.Sprintf("posts:%v:page:%d", version, pageNum)
	err = p.rdb.ZAdd(ctx, valueKey, postsKeys...).Err()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	err = p.rdb.Expire(ctx, valueKey, 1*time.Minute).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	err = p.rdb.Expire(ctx, metaKey, 1*time.Minute).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	return nil
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
	return nil
}

func (p *postRedis) PinPost(ctx context.Context, id string) error {
	key := fmt.Sprintf("posts:%s", id)

	pinnedStr, err := p.rdb.HGet(ctx, key, "pinned").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	pinned, _ := strconv.ParseBool(pinnedStr)
	pinned = !pinned
	_, err = p.rdb.HSet(ctx, key, "pinned", pinned).Result()
	if err != nil {
		return errors.Join(err, ErrInternalFailure)
	}
	return nil
}

func (p *postRedis) Update(ctx context.Context, id string, post types.Post) error {
	panic("unimplemented")
}
