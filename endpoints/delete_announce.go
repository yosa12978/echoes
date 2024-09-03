package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func DeleteAnnounce(logger logging.Logger, service services.Announce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.Delete(r.Context()); err != nil {
			utils.RenderBlock(w, "alert", "failed to delete announce")
			return
		}
		logger.Info("announce removed")
		utils.RenderBlock(w, "alert", "Announce deleted")
	}
}
