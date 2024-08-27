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

type Comment interface {
	GetComment(ctx context.Context) http.Handler
	GetPostComments(ctx context.Context) http.Handler
	GetCommentCount(ctx context.Context) http.Handler
	CreateComment(ctx context.Context) http.Handler
	DeleteComment(ctx context.Context) http.Handler
}

type comment struct {
	commentService services.Comment
	logger         logging.Logger
}

func NewComment(commentService services.Comment, logger logging.Logger) Comment {
	return &comment{commentService: commentService, logger: logger}
}

func (h *comment) GetComment(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		comment, err := h.commentService.GetCommentById(ctx, mux.Vars(r)["id"])
		if err != nil {
			utils.RenderBlock(w, "alert", "comment not found")
			return
		}
		utils.RenderBlock(w, "comment", comment)
	})
}

func (h *comment) GetPostComments(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := r.URL.Query().Get("postId") // temporary solution
		pagestr := r.URL.Query().Get("page")
		if pagestr == "" {
			pagestr = "1"
		}
		page, err := strconv.Atoi(pagestr)
		if err != nil {
			utils.RenderBlock(w, "alert", "wrong page number")
			return
		}
		commentsPaged, err := h.commentService.GetPostComments(ctx, postId, page, 20)
		if err != nil {
			h.logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't fetch post comments")
			return
		}
		utils.RenderBlock(w, "comments", commentsPaged)
	})
}

func (h *comment) CreateComment(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := r.URL.Query().Get("postId")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		_, err := h.commentService.CreateComment(
			ctx,
			postId,
			body["name"].(string),
			body["email"].(string),
			body["content"].(string),
		)
		if err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		utils.RenderBlock(w, "alert", "comment created")
	})
}

func (h *comment) DeleteComment(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commentId := mux.Vars(r)["id"]
		_, err := h.commentService.DeleteComment(ctx, commentId)
		if err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		utils.RenderBlock(w, "alert", "comment deleted")
	})
}

func (h *comment) GetCommentCount(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, _ := h.commentService.GetCommentsCount(ctx, mux.Vars(r)["id"])
		w.Write([]byte(strconv.Itoa(count)))
	})
}
