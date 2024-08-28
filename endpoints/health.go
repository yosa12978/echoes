package endpoints

import (
	"net/http"

	"github.com/yosa12978/echoes/services"
)

func Healthcheck(service services.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := service.Healthcheck(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	}
}
