package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func CreateLink(logger logging.Logger, service services.Link) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, problems, err := utils.
			ReadJsonAndValidate[types.LinkCreateDto](r.Context(), r.Body)
		if err != nil {
			if errors.Is(err, types.ErrValidationFailed) {
				utils.RenderBlock(w, "problems", problems)
				return
			}
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}
		place, err := strconv.Atoi(dto.Place)
		if err != nil {
			utils.RenderBlock(w, "alert_danger", "place must be a number")
			return
		}
		_, err = service.CreateLink(r.Context(), dto.Name, dto.URL, dto.Icon, place)
		if err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}
		utils.RenderBlock(w, "alert_success", "Created new link")
	}
}
