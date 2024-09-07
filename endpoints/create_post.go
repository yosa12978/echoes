package endpoints

import (
	"errors"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func CreatePost(logger logging.Logger, service services.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, problems, err := utils.
			ReadJsonAndValidate[types.PostCreateDto](r.Context(), r.Body)
		if err != nil {
			if errors.Is(err, types.ErrValidationFailed) {
				utils.RenderBlock(w, "problems", problems)
				return
			}
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}
		if _, err := service.CreatePost(
			r.Context(),
			dto.Title,
			dto.Content,
			dto.Tweet != "",
		); err != nil {
			logger.Error(err.Error())
			utils.RenderBlock(w, "alert_danger", "Failed to create")
			return
		}
		utils.RenderBlock(w, "alert_success", "Created new post")
	}
}
