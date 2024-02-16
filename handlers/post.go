package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Post interface {
	GetPosts(ctx context.Context) http.Handler
	GetPostById(ctx context.Context) http.Handler
	CreatePost(ctx context.Context) http.Handler
	DeletePost(ctx context.Context) http.Handler
}

type post struct {
	postRepo repos.Post
}

func NewPost(postRepo repos.Post) Post {
	return &post{
		postRepo: postRepo,
	}
}

func (h *post) GetPosts(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		pageS := r.URL.Query().Get("page")
		if pageS == "" {
			pageS = "1"
		}
		page, err := strconv.Atoi(pageS)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		limitS := r.URL.Query().Get("limit")
		if limitS == "" {
			limitS = "5"
		}
		limit, err := strconv.Atoi(limitS)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		res := types.Payload{
			Title:   "Home",
			Content: h.postRepo.GetPage(ctx, page, limit),
		}
		utils.RenderBlock(w, "postsPage", res)
	})
}

func (h *post) GetPostById(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post, err := h.postRepo.FindById(ctx, id)
		if err != nil {
			log.Println(err.Error())
			return
		}
		utils.RenderBlock(w, "detail", post)
	})
}

func (h *post) CreatePost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// validation!!!
		post := types.NewPost(r.FormValue("title"), r.FormValue("content"))
		if _, err := h.postRepo.Create(ctx, post); err != nil {
			log.Println(err.Error())
			return
		}
		utils.RenderBlock(w, "alert", "Created new post")
	})
}

func (h *post) DeletePost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, err := h.postRepo.Delete(ctx, id); err != nil {
			log.Println(err.Error())
			return
		}
	})
}
