package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func DeletePost(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if _, err := service.DeletePost(r.Context(), r.FormValue("id")); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert_danger", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert_success", "Post deleted")
	}
}
