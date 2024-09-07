package endpoints

import (
	"errors"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func CreateAnnounce(logger logging.Logger, service services.Announce) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, problems, err := utils.
			ReadJsonAndValidate[types.AnnounceCreateDto](r.Context(), r.Body)
		if err != nil {
			if errors.Is(err, types.ErrValidationFailed) {
				utils.RenderBlock(w, "problems", problems)
				return
			}
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}

		if err := service.Create(r.Context(), dto.Content); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert_danger", "can't create announce")
			return
		}
		utils.RenderBlock(w, "alert_success", "Announce created")
	}
}
