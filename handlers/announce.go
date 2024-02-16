package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

type Announce interface {
	Get(ctx context.Context) http.Handler
}

type announce struct {
	announceRepo repos.Announce
}

func NewAnnounce(announceRepo repos.Announce) Announce {
	handler := new(announce)
	handler.announceRepo = announceRepo
	return handler
}

func (h *announce) Get(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.RenderBlock(w, "announce", h.announceRepo.Get())
	})
}
