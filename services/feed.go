package services

import (
	"bytes"
	"context"
	"time"

	"github.com/gorilla/feeds"
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
	feed := &feeds.Feed{
		Title:       "yusuf's microblog recent posts",
		Link:        &feeds.Link{Href: "https://blinkk.org/blog"},
		Description: "30 latest posts from yusuf's microblog",
		Author:      &feeds.Author{Name: "Yusuf Yakubov", Email: "yosa12978@gmail.com"},
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
			Link:    &feeds.Link{Href: "https://blinkk.org/posts/" + v.Id},
			Author:  &feeds.Author{Name: "Yusuf Yakubov", Email: "yosa12978@gmail.com"},
			Created: created,
		}
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(v.Content), &buf); err != nil {
			return err.Error(), err
		}
		item.Content = buf.String()
		items = append(items, item)
	}
	feed.Items = items
	return feed.ToAtom()
}
