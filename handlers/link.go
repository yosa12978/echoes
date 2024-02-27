package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Link interface {
	GetAll(ctx context.Context) http.Handler
	GetAdmin(ctx context.Context) http.Handler
	Create(ctx context.Context) http.Handler
	Delete(ctx context.Context) http.Handler
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
		links, err := h.linkRepo.FindAll(ctx)
		if err != nil {
			// log here
		}
		utils.RenderBlock(w, "links", links)
	})
}

func (h *link) GetAdmin(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		links, err := h.linkRepo.FindAll(ctx)
		if err != nil {
			// log here
		}
		utils.RenderBlock(w, "links_admin", links)
	})
}

func (h *link) Create(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("name")
		url := r.FormValue("url")
		link := types.Link{
			Name: name,
			URL:  url,
		}
		_, err := h.linkRepo.Create(ctx, link)
		if err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Failed to create")
			return
		}
		utils.RenderBlock(w, "alert", "Created new link")
	})
}

func (h *link) Delete(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, err := h.linkRepo.Delete(ctx, id); err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert", "Link Deleted")
	})
}
