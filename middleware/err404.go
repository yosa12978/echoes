package middleware

import (
	"net/http"

	"github.com/yosa12978/echoes/utils"
)

func Err404() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if r.Pattern == "" { // shitty idea. I need to completely flush http.ResponseWriter for it to work
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(404)
				if err := utils.RenderView(w, "err404", "Error 404", nil); err != nil {
					http.Error(w, err.Error(), 500)
				}
			}
		})
	}
}
