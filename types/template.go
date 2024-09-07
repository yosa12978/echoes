package types

type CommentsInfo struct {
	Page[Comment]
	PostId string
}

type Templ struct {
	Title   string
	Logo    string
	BgImg   string
	Payload interface{}
}
