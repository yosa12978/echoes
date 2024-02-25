package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

type Comment interface {
	GetPostComments(ctx context.Context) http.Handler
	CreateComment(ctx context.Context) http.Handler
	DeleteComment(ctx context.Context) http.Handler
}

type comment struct {
	commentService services.Comment
}

func NewComment(commentService services.Comment) Comment {
	return &comment{commentService: commentService}
}

func (h *comment) GetPostComments(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := r.URL.Query().Get("postId") // temporary solution
		pagestr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pagestr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		commentsPaged, err := h.commentService.GetPostComments(ctx, postId, page, 20)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		utils.RenderBlock(w, "comments", commentsPaged)
	})
}

func (h *comment) CreateComment(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		postId := r.URL.Query().Get("postId")
		name := r.FormValue("name")
		email := r.FormValue("email")
		content := r.FormValue("content")
		_, err := h.commentService.CreateComment(ctx, postId, name, email, content)
		if err != nil {
			http.Error(w, err.Error(), 400)
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		utils.RenderBlock(w, "alert", "comment created")
	})
}

func (h *comment) DeleteComment(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this handler is not implemented yet"))
	})
}
