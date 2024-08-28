package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetLinks(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := service.GetLinks(r.Context())
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "links", links)
			return
		}
		utils.RenderBlock(w, "links", links)
	}
}
