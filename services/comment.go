package services

import "github.com/yosa12978/echoes/repos"

type Comment interface {
}

type comment struct {
	commentRepo repos.Comment
}

func NewComment(commentRepo repos.Comment) Comment {
	return &comment{commentRepo: commentRepo}
}
