package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Post interface {
	GetPosts(ctx context.Context) http.Handler
	GetPostById(ctx context.Context) http.Handler
	PinPost(ctx context.Context) http.Handler
	CreatePost(ctx context.Context) http.Handler
	DeletePost(ctx context.Context) http.Handler
	Search(ctx context.Context) http.Handler
}

type post struct {
	postService services.Post
	logger      logging.Logger
}

func NewPost(postService services.Post, logger logging.Logger) Post {
	return &post{
		postService: postService,
		logger:      logger,
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
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "wrond page number")
			return
		}
		limitS := r.URL.Query().Get("limit")
		if limitS == "" {
			limitS = "20"
		}
		limit, err := strconv.Atoi(limitS)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "wrond limit number")
			return
		}
		searchQuery := r.URL.Query().Get("query")
		var posts *types.Page[types.Post]
		if searchQuery == "" {
			posts, err = h.postService.GetPostsPaged(ctx, page, limit)
		} else {
			posts, err = h.postService.Search(ctx, searchQuery, page, limit)
		}
		if err != nil {
			h.logger.Error(err)
		}
		if posts.Total == 0 {
			utils.RenderBlock(w, "noPosts", nil)
			return
		}
		utils.RenderBlock(w, "postsPage", posts)
	})
}

func (h *post) GetPostById(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post, err := h.postService.GetPostById(ctx, id) //h.postRepo.FindById(ctx, id)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "post not found")
			return
		}
		utils.RenderBlock(w, "post", post)
	})
}

func (h *post) CreatePost(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		title := r.FormValue("title")
		content := r.FormValue("content")
		tweetCheckbox := r.FormValue("tweet")
		tweet := false
		if tweetCheckbox == "on" {
			tweet = true
		}
		fmt.Printf("tweetCheckbox: %v\n", tweetCheckbox)
		fmt.Printf("tweet: %v\n", tweet)
		if _, err := h.postService.CreatePost(ctx, title, content, tweet); err != nil {
			h.logger.Error(err)
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
		if _, err := h.postService.DeletePost(ctx, body["id"].(string)); err != nil {
			h.logger.Error(err)
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
		_, err := h.postService.PinPost(ctx, body["id"].(string))
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "Post not found")
			return
		}
		utils.RenderBlock(w, "alert", "Post pinned :)")
	})
}

// depricated
func (h *post) Search(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		pageS := r.URL.Query().Get("page")
		if pageS == "" {
			pageS = "1"
		}
		page, err := strconv.Atoi(pageS)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "wrond page number")
			return
		}
		limitS := r.URL.Query().Get("limit")
		if limitS == "" {
			limitS = "20"
		}
		limit, err := strconv.Atoi(limitS)
		if err != nil {
			h.logger.Error(err)
			utils.RenderBlock(w, "alert", "wrond limit number")
			return
		}
		posts, err := h.postService.Search(ctx, q, page, limit)
		if err != nil {
			h.logger.Error(err)
		}
		if posts.Total == 0 {
			utils.RenderBlock(w, "noPosts", nil)
			return
		}
		utils.RenderBlock(w, "postsPage", posts)
	})
}
