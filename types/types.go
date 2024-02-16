package types

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id      string
	Title   string
	Content string
	Created string
}

func NewPost(title, content string) Post {
	return Post{
		Id:      uuid.NewString(),
		Title:   title,
		Content: content,
		Created: time.Now().Format(time.RFC3339),
	}
}

type Link struct {
	Id      string
	Name    string
	URL     string
	Created string
}

type Comment struct {
	Id      string
	Email   string
	Name    string
	Content string
	Created string
}

type Account struct {
	Id       string
	Username string
	Password string
	Created  string
}

type Announce struct {
	Content string
}

type Config struct {
	Addr     string
	Postgres string
}

type Page[T interface{}] struct {
	HasNext  bool
	Size     int
	NextPage int
	Content  []T
	Total    int
}

type Payload struct {
	Title   string
	Content interface{}
}
