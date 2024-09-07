package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func PinPost(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		_, err := service.PinPost(r.Context(), body["id"].(string))
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert_danger", "Post not found")
			return
		}
		utils.RenderBlock(w, "alert_success", "Post pinned")
	}
}
