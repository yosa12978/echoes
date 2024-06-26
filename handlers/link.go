package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	Portal(ctx context.Context) http.Handler
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
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		name := body["name"].(string)
		url := body["url"].(string)
		icon := body["icon"].(string)
		placeStr := body["place"].(string)
		place, err := strconv.Atoi(placeStr)
		if err != nil {
			utils.RenderBlock(w, "alert", "place must be a number")
			return
		}
		_, err = h.linkService.CreateLink(ctx, name, url, icon, place)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", err.Error())
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

func (h *link) Portal(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if id == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		link, err := h.linkService.GetLinkById(ctx, id)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "not found")
			return
		}
		http.Redirect(w, r, link.URL, http.StatusMovedPermanently)
	})
}
