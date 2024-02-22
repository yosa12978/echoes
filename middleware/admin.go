package middleware

import (
	"net/http"

	"github.com/yosa12978/echoes/session"
)

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionInfo, err := session.GetInfo(r)
		if err != nil || sessionInfo == nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if sessionInfo.Role != "ADMIN" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
