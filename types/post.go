package types

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id       string
	Title    string
	Content  string
	Created  string
	Pinned   bool
	Tweet    bool
	Comments int
}

func NewPost(title, content string, tweet bool) Post {
	return Post{
		Id:       uuid.NewString(),
		Title:    title,
		Content:  content,
		Created:  time.Now().Format(time.RFC3339),
		Pinned:   false,
		Tweet:    tweet,
		Comments: 0,
	}
}

type PostCreateDto struct {
	Title   string
	Content string
	Tweet   string
}

func (p *PostCreateDto) Validate(ctx context.Context) (problems map[string]string, ok bool) {
	problems = make(map[string]string)
	p.Content = strings.TrimSpace(p.Content)
	p.Title = strings.TrimSpace(p.Title)
	if p.Content == "" {
		problems["content"] = "content can't be empty"
	}
	if p.Title == "" {
		problems["title"] = "title can't be empty"
	}
	return problems, len(problems) == 0
}

type PostUpdateDto struct {
	Title   string
	Content string
	Tweet   string
}

func (p *PostUpdateDto) Validate(ctx context.Context) (problems map[string]string, ok bool) {
	problems = make(map[string]string)
	p.Content = strings.TrimSpace(p.Content)
	p.Title = strings.TrimSpace(p.Title)
	if p.Content == "" {
		problems["content"] = "content can't be empty"
	}
	if p.Title == "" {
		problems["title"] = "title can't be empty"
	}
	return problems, len(problems) == 0
}
