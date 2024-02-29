package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

type Announce interface {
	Get(ctx context.Context) http.Handler
	Create(ctx context.Context) http.Handler
	Delete(ctx context.Context) http.Handler
}

type announce struct {
	announceService services.Announce
	logger          logging.Logger
}

func NewAnnounce(announceService services.Announce, logger logging.Logger) Announce {
	handler := new(announce)
	handler.announceService = announceService
	handler.logger = logger
	return handler
}

func (h *announce) Get(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		announce, err := h.announceService.Get(ctx)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "can't fetch announce")
			return
		}
		utils.RenderBlock(w, "announce", announce)
	})
}

func (h *announce) Create(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		_, err := h.announceService.Create(ctx, r.FormValue("content"))
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "can't create announce")
			return
		}
		utils.RenderBlock(w, "alert", "Announce created")
	})
}

func (h *announce) Delete(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.announceService.Delete(ctx); err != nil {
			utils.RenderBlock(w, "alert", "failed to delete announce")
			return
		}
		utils.RenderBlock(w, "alert", "Announce deleted")
	})
}
