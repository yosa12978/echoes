package endpoints

import (
	"net/http"
)

func Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// i don't need this one at the time
	}
}
