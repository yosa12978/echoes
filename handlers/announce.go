package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

type Announce interface {
	Get(ctx context.Context) http.Handler
	Create(ctx context.Context) http.Handler
	Delete(ctx context.Context) http.Handler
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

func (h *announce) Create(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			utils.RenderBlock(w, "alert", "Can't create an empty announce")
			return
		}
		h.announceRepo.Create(content)
		utils.RenderBlock(w, "alert", "Announce created")
	})
}

func (h *announce) Delete(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.announceRepo.Delete()
		utils.RenderBlock(w, "alert", "Announce deleted")
	})
}
