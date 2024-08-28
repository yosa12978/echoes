package endpoints

import (
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func GetPosts(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageS := r.URL.Query().Get("page")
		if pageS == "" {
			pageS = "1"
		}
		page, err := strconv.Atoi(pageS)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "wrond page number")
			return
		}
		limitS := r.URL.Query().Get("limit")
		if limitS == "" {
			limitS = "20"
		}
		limit, err := strconv.Atoi(limitS)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "wrond limit number")
			return
		}
		searchQuery := r.URL.Query().Get("query")
		var posts *types.Page[types.Post]
		if searchQuery == "" {
			posts, err = service.GetPostsPaged(r.Context(), page, limit)
		} else {
			posts, err = service.Search(r.Context(), searchQuery, page, limit)
		}
		if err != nil {
			logger.Error(err.Error())
		}
		if posts.Total == 0 {
			utils.RenderBlock(w, "noPosts", nil)
			return
		}
		utils.RenderBlock(w, "postsPage", posts)
	}
}
