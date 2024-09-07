package types

import (
	"context"
	"strings"
)

type Announce struct {
	Content string
	Date    string
}

type AnnounceCreateDto struct {
	Content string
}

func (a AnnounceCreateDto) Validate(ctx context.Context) (
	AnnounceCreateDto,
	map[string]string,
	bool,
) {
	problems := make(map[string]string)
	a.Content = strings.TrimSpace(a.Content)
	if a.Content == "" {
		problems["content"] = "content can't be empty"
	}
	return a, problems, len(problems) == 0
}
