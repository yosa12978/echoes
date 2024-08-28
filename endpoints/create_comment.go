package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func CreateComment(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId := r.URL.Query().Get("postId")
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		_, err := service.CreateComment(
			r.Context(),
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
	}
}
