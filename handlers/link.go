package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

type Link interface {
	GetAll(ctx context.Context) http.Handler
	GetAdmin(ctx context.Context) http.Handler
	Create(ctx context.Context) http.Handler
	Delete(ctx context.Context) http.Handler
}

type link struct {
	logger      logging.Logger
	linkService services.Link
}

func NewLink(linkService services.Link, logger logging.Logger) Link {
	handler := new(link)
	handler.linkService = linkService
	handler.logger = logger
	return handler
}

func (h *link) GetAll(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		links, err := h.linkService.GetLinks(ctx)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "links", links)
			return
		}
		utils.RenderBlock(w, "links", links)
	})
}

// merge this method with main one
func (h *link) GetAdmin(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		links, err := h.linkService.GetLinks(ctx)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "can't fetch links")
			return
		}
		utils.RenderBlock(w, "links_admin", links)
	})
}

func (h *link) Create(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("name")
		url := r.FormValue("url")
		_, err := h.linkService.CreateLink(ctx, name, url)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "Failed to create")
			return
		}
		utils.RenderBlock(w, "alert", "Created new link")
	})
}

func (h *link) Delete(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, err := h.linkService.DeleteLink(ctx, id); err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert", "Link Deleted")
	})
}
