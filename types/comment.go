package types

import (
	"context"
	"net/mail"
	"strings"
)

type Comment struct {
	Id      string
	Email   string
	Name    string
	Content string
	Created string
	PostId  string
}

type CommentCreateDto struct {
	Name    string
	Email   string
	Content string
}

func (c CommentCreateDto) Validate(ctx context.Context) (CommentCreateDto, map[string]string, bool) {
	problems := make(map[string]string)
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	c.Content = strings.TrimSpace(c.Content)

	if c.Name == "" {
		problems["name"] = "Name is required"
	}
	if c.Email == "" {
		problems["email"] = "Email is required"
	}
	if c.Content == "" {
		problems["content"] = "Content can't be empty"
	}
	if _, err := mail.ParseAddress(c.Email); err != nil {
		problems["email"] = "Email is invalid"
	}

	return c, problems, len(problems) == 0
}
