package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func Portal(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		link, err := service.GetLinkById(r.Context(), id)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "not found")
			return
		}
		http.Redirect(w, r, link.URL, http.StatusMovedPermanently)
	}
}
