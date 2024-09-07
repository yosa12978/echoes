package types

import (
	"context"
	"net/url"
	"strings"
)

type Link struct {
	Id      string
	Name    string
	URL     string
	Created string
	Icon    string
	Place   int
}

type LinkCreateDto struct {
	Name    string
	URL     string
	Created string
	Icon    string
	Place   string // because html form interpretes numeric form as string
}

func (l *LinkCreateDto) Validate(ctx context.Context) (problems map[string]string, ok bool) {
	problems = make(map[string]string)
	l.Name = strings.TrimSpace(l.Name)
	l.URL = strings.TrimSpace(l.URL)
	if l.Name == "" {
		problems["name"] = "name can't be empty"
	}
	_, err := url.ParseRequestURI(l.URL)
	if err != nil {
		problems["url"] = "url must be valid"
	}
	if l.Place == "" {
		l.Place = "1"
	}
	return problems, len(problems) == 0
}

type LinkUpdateDto struct {
	Name  string
	URL   string
	Icon  string
	Place string
}

func (l *LinkUpdateDto) Validate(ctx context.Context) (problems map[string]string, ok bool) {
	problems = make(map[string]string)
	l.Name = strings.TrimSpace(l.Name)
	l.URL = strings.TrimSpace(l.URL)
	if strings.TrimSpace(l.Name) == "" {
		problems["name"] = "name can't be empty"
	}
	_, err := url.ParseRequestURI(l.URL)
	if err != nil {
		problems["url"] = "url must be valid"
	}
	if l.Place == "" {
		l.Place = "1"
	}
	return problems, len(problems) == 0
}
