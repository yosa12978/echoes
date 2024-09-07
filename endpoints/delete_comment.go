package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func DeleteComment(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		commentId := r.FormValue("id")
		_, err := service.DeleteComment(r.Context(), commentId)
		if err != nil {
			utils.RenderBlock(w, "alert_danger", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert_success", "Comment deleted")
	}
}
