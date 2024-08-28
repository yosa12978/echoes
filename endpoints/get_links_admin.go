package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

// merge this method with main one
func GetLinksAdmin(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := service.GetLinks(r.Context())
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert", "can't fetch links")
			return
		}
		utils.RenderBlock(w, "links_admin", links)
	}
}
