package services

import "github.com/yosa12978/echoes/repos"

type Post interface {
}

type post struct {
	postRepo repos.Post
}

func NewPost(postRepo repos.Post) Post {
	return &post{postRepo: postRepo}
}
