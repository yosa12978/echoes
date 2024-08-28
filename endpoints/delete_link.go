package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func DeleteLink(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, err := service.DeleteLink(r.Context(), id); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "Failed to delete")
			return
		}
		utils.RenderBlock(w, "alert", "Link Deleted")
	}
}
