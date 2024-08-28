package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetPostById(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		post, err := service.GetPostById(r.Context(), id)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "post not found")
			return
		}
		utils.RenderBlock(w, "post", post)
	}
}
