package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func CreateLink(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		name := body["name"].(string)
		url := body["url"].(string)
		icon := body["icon"].(string)
		placeStr := body["place"].(string)
		place, err := strconv.Atoi(placeStr)
		if err != nil {
			utils.RenderBlock(w, "alert", "place must be a number")
			return
		}
		_, err = service.CreateLink(r.Context(), name, url, icon, place)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		utils.RenderBlock(w, "alert", "Created new link")
	}
}
