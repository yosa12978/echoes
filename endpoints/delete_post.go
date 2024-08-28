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
		// var body map[string]interface{}
		// if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		// 	h.logger.Error(err.Error())
		// 	return
		// }
		if _, err := service.DeletePost(r.Context(), r.FormValue("id")); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert", "Post deleted successfully")
	}
}
