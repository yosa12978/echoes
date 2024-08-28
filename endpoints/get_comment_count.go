package endpoints

import (
	"net/http"
	"strconv"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
)

func GetCommentCount(logger logging.Logger, service services.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, _ := service.GetCommentsCount(r.Context(), r.PathValue("id"))
		w.Write([]byte(strconv.Itoa(count)))
	}
}
