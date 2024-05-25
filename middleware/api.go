package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func ApiWrapper(next types.ApiFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiResp, err := next(w, r)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(apiResp.Code)
		if err != nil {
			switch r.Header.Get("Content-Type") {
			case "application/json":
				json.NewEncoder(w).Encode(
					types.ApiMsg{
						Code:    apiResp.Code,
						Message: err.Error(),
					},
				)
				return
			case "text/html":
				utils.RenderBlock(w, "alert", err.Error())
				return
			}
		}
		if r.Header.Get("Content-Type") == "application/json" {
			json.NewEncoder(w).Encode(apiResp.Body)
			return
		}
		utils.RenderBlock(w, apiResp.Templ, apiResp.Body)
	})
}
