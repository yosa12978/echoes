package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/session"
)

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := session.EndSession(r, w)
		if err != nil {
			http.Error(w, "you can't logout unless you logged in", http.StatusUnauthorized)
			return
		}
		w.Header().Set("HX-Redirect", "/")
	}
}
