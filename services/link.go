package services

import "github.com/yosa12978/echoes/repos"

type Link interface {
}

type link struct {
	linkRepo repos.Link
}

func NewLink(linkRepo repos.Link) Link {
	return &link{linkRepo: linkRepo}
}
