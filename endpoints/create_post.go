package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func CreatePost(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		title := body["title"].(string)
		content := body["content"].(string)
		_, tweet := body["tweet"]

		if _, err := service.CreatePost(
			r.Context(),
			title,
			content,
			tweet,
		); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "Failed to create")
			return
		}
		utils.RenderBlock(w, "alert", "Created new post")
	}
}
