package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

type Link interface {
	GetAll(ctx context.Context) http.Handler
}

type link struct {
	linkRepo repos.Link
}

func NewLink(linkRepo repos.Link) Link {
	handler := new(link)
	handler.linkRepo = linkRepo
	return handler
}

func (h *link) GetAll(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RenderBlock(w, "links", h.linkRepo.FindAll(ctx))
	})
}
