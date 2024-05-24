package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func ApiWrapper(next types.ApiFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, code, err := next(w, r)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(code)
		if err != nil {
			switch r.Header.Get("Content-Type") {
			case "application/json":
				json.NewEncoder(w).Encode(
					types.ApiError{
						StatusCode: code,
						Err:        err.Error(),
					},
				)
				return
			case "text/html":
				utils.RenderBlock(w, "alert", err.Error())
				return
			}
		}
		if r.Header.Get("Content-Type") == "application/json" {
			json.NewEncoder(w).Encode(body)
			return
		}
		// render requested template
	})
}
