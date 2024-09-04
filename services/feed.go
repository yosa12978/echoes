package services

import (
	"bytes"
	"context"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yosa12978/echoes/config"
	"github.com/yosa12978/echoes/types"
	"github.com/yuin/goldmark"
)

type Feed interface {
	GenerateFeed(ctx context.Context) (string, error)
}

type feed struct {
	postService Post
}

func NewFeedService(postService Post) Feed {
	return &feed{
		postService: postService,
	}
}

// TODO: refactor
func (f *feed) GenerateFeed(ctx context.Context) (string, error) {
	cfg := config.Get()
	feed := &feeds.Feed{
		Title:       cfg.Feed.Title,
		Link:        &feeds.Link{Href: cfg.Feed.Link},
		Description: cfg.Feed.Desc,
		Author:      &feeds.Author{Name: cfg.Feed.Author, Email: cfg.Feed.Email},
		Created:     time.Now().UTC(),
	}
	items := []*feeds.Item{}
	posts, err := f.postService.GetPostsPaged(ctx, 1, 30)
	if err != nil {
		return "couldn't get posts", err
	}
	for _, v := range posts.Content {
		created, _ := time.Parse(time.RFC3339, v.Created)
		item := &feeds.Item{
			Id:      v.Id,
			Title:   v.Title,
			Link:    &feeds.Link{Href: cfg.Feed.DetailLink + v.Id},
			Author:  &feeds.Author{Name: cfg.Feed.Author, Email: cfg.Feed.Email},
			Created: created,
		}
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(v.Content), &buf); err != nil {
			return err.Error(), types.NewErrInternalFailure(err)
		}
		item.Content = buf.String()
		items = append(items, item)
	}
	feed.Items = items
	return feed.ToAtom()
}
