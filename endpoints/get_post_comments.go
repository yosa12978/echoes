package endpoints

import (
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetPostComments(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		commentsPaged, err := service.GetPostComments(r.Context(), postId, page, 20)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't fetch post comments")
			return
		}
		utils.RenderBlock(w, "comments", commentsPaged)
	}
}
