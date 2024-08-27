package router

import (
	"net/http"

	"github.com/yosa12978/echoes/middleware"
)

func New(opts ...optionFunc) http.Handler {
	options := newOptions(opts...)
	router := http.NewServeMux()
	addRoutes(router, options)
	var handler http.Handler = router
	handler = middleware.Pipeline(
		router,
		middleware.Latency(options.logger),
		middleware.StripSlash,
		middleware.Recovery(options.logger),
	)
	return handler
}

func addRoutes(r *http.ServeMux, options options) {

}
