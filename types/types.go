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
	Pinned  bool
}

func NewPost(title, content string) Post {
	return Post{
		Id:      uuid.NewString(),
		Title:   title,
		Content: content,
		Created: time.Now().Format(time.RFC3339),
		Pinned:  false,
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
	PostId  string
}

type Account struct {
	Id       string
	Username string
	Password string
	Created  string
	IsAdmin  bool
}

type Announce struct {
	Content string
	Date    string
}

type Config struct {
	Addr       string
	Postgres   string
	SessionKey string
	RedisAddr  string
	RedisDb    int
	RedisPwd   string
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

type Profile struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
	Icon string `json:"icon"`
}

type SessionInfo struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}

type CommentsInfo struct {
	Page[Comment]
	PostId string
}
