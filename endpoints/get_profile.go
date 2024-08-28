package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetProfile(logger logging.Logger, service services.Profile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := service.Get(r.Context())
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't render profile information")
			return
		}
		utils.RenderBlock(w, "profile", res)
	}
}
