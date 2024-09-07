package endpoints

import (
	"errors"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func CreateComment(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto, problems, err := utils.
			ReadJsonAndValidate[types.CommentCreateDto](r.Context(), r.Body)
		if err != nil {
			if errors.Is(err, types.ErrValidationFailed) {
				utils.RenderBlock(w, "problems", problems)
				return
			}
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}

		postId := r.URL.Query().Get("postId")

		_, err = service.CreateComment(
			r.Context(),
			postId,
			dto.Name,
			dto.Email,
			dto.Content,
		)
		if err != nil {
			utils.RenderBlock(w, "alert_danger", err.Error())
			return
		}
		utils.RenderBlock(w, "alert_success", "comment created")
	}
}
