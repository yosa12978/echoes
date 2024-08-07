package types

import (
	"net/http"
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

type Link struct {
	Id      string
	Name    string
	URL     string
	Created string
	Icon    string
	Place   int
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
	Salt     string
	Created  string
	IsAdmin  bool
}

type Announce struct {
	Content string
	Date    string
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

type Session struct {
	Username        string `json:"username"`
	IsAdmin         bool   `json:"is_admin"`
	Timestamp       int64  `json:"timestamp"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

type CommentsInfo struct {
	Page[Comment]
	PostId string
}

type Templ struct {
	Title   string
	Payload interface{}
}

// Backlog: for future api or smth
type ApiFunc func(w http.ResponseWriter, r *http.Request) (ApiResponse, error)

type ApiResponse struct {
	Body  interface{}
	Templ string // for htmx
	Code  int
}

type ApiMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewApiResp(
	Body interface{},
	Templ string,
	Code int,
) ApiResponse {
	return ApiResponse{
		Body:  Body,
		Templ: Templ,
		Code:  Code,
	}
}
