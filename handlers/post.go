package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Post interface {
	GetPosts(ctx context.Context) http.Handler
	GetPostById(ctx context.Context) http.Handler
	PinPost(ctx context.Context) http.Handler
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
			limitS = "20"
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
		if res.Content.(*types.Page[types.Post]).Total == 0 {
			utils.RenderBlock(w, "noPosts", nil)
			return
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
			http.Error(w, err.Error(), 404)
			return
		}
		pl := types.Payload{Title: post.Title, Content: post}
		utils.RenderBlock(w, "post", pl)
	})
}

func (h *post) CreatePost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// validation!!!
		post := types.NewPost(r.FormValue("title"), r.FormValue("content"))
		if _, err := h.postRepo.Create(ctx, post); err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Failed to create")
			return
		}
		utils.RenderBlock(w, "alert", "Created new post")
	})
}

func (h *post) DeletePost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if _, err := h.postRepo.Delete(ctx, body["id"].(string)); err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert", "Post deleted successfully")
	})
}

func (h *post) PinPost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		post, err := h.postRepo.FindById(ctx, body["id"].(string))
		if err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Post not found")
			return
		}
		post.Pinned = !post.Pinned
		_, err = h.postRepo.Update(ctx, post.Id, *post)
		if err != nil {
			log.Println(err.Error())
			utils.RenderBlock(w, "alert", "Failed to update")
			return
		}
		utils.RenderBlock(w, "alert", "Post pinned :)")
	})
}
