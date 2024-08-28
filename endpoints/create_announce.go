package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func CreateAnnounce(logger logging.Logger, service services.Announce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		content := body["content"].(string)

		// r.ParseForm()
		// content := r.FormValue("content")
		_, err := service.Create(r.Context(), content)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't create announce")
			return
		}
		utils.RenderBlock(w, "alert", "Announce created")
	}
}
