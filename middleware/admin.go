package middleware

import (
	"net/http"

	"github.com/yosa12978/echoes/session"
)

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := session.GetSession(r)
		if err != nil || s == nil {
			http.Error(w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		if !s.IsAuthenticated {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !s.IsAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
