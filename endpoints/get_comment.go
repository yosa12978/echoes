package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/utils"
)

func GetComment(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comment, err := service.GetCommentById(r.Context(), r.PathValue("id"))
		if err != nil {
			utils.RenderBlock(w, "alert", "comment not found")
			return
		}
		utils.RenderBlock(w, "comment", comment)
	}
}
