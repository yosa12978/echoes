package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/services"
)

func GetFeed(service services.Feed) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		feed, err := service.GenerateFeed(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(feed))
	}
}
