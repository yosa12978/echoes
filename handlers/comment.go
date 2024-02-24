package handlers

type Comment interface {
}
type comment struct {
}

func NewComment() Comment {
	return &comment{}
}
