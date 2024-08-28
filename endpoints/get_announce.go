package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetAnnounce(logger logging.Logger, service services.Announce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		announce, err := service.Get(r.Context())
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't fetch announce")
			return
		}
		utils.RenderBlock(w, "announce", announce)
	}
}
