package router

import (
	"net/http"

	"github.com/yosa12978/echoes/endpoints"
	"github.com/yosa12978/echoes/middleware"
	"github.com/yosa12978/echoes/session"
	"github.com/yosa12978/echoes/utils"
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
	apiRouter := http.NewServeMux()

	addLinkRoutes(apiRouter, options)
	addPostRoutes(apiRouter, options)
	addProfileRoutes(apiRouter, options)
	addFeedRoutes(r, options)
	addAccountRoutes(apiRouter, options)
	addAnnounceRoutes(apiRouter, options)
	addHealthRoutes(r, options)
	addCommentRoutes(apiRouter, options)
	addViewRoutes(r)

	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	r.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("./assets/")),
	))
}

func addLinkRoutes(router *http.ServeMux, options options) {
	router.Handle("GET /links",
		endpoints.GetLinks(options.logger, options.linkService))

	router.Handle("GET /portal/{id}",
		endpoints.Portal(options.logger, options.linkService))

	router.Handle("GET /links-admin", middleware.Admin(
		endpoints.GetLinksAdmin(options.logger, options.linkService),
	))

	router.Handle("POST /links", middleware.Admin(
		endpoints.CreateLink(options.logger, options.linkService),
	))

	router.Handle("DELETE /links/{id}",
		middleware.Admin(
			endpoints.DeleteLink(options.logger, options.linkService),
		),
	)
}

func addPostRoutes(router *http.ServeMux, options options) {
	router.Handle("GET /posts",
		endpoints.GetPosts(options.logger, options.postService))

	router.Handle("GET /posts/{id}",
		endpoints.GetPostById(options.logger, options.postService))

	router.Handle("POST /posts",
		middleware.Admin(
			endpoints.CreatePost(options.logger, options.postService),
		),
	)

	router.Handle("DELETE /posts",
		middleware.Admin(
			endpoints.DeletePost(options.logger, options.postService),
		),
	)

	router.Handle("PATCH /post-pin",
		middleware.Admin(
			endpoints.PinPost(options.logger, options.postService),
		),
	)
}

func addCommentRoutes(router *http.ServeMux, options options) {
	router.Handle("GET /comments",
		endpoints.GetPostComments(options.logger, options.commentService))

	router.Handle("POST /comments",
		endpoints.CreateComment(options.logger, options.commentService))

	router.Handle("DELETE /comments/{id}",
		middleware.Admin(
			endpoints.DeleteComment(options.logger, options.commentService),
		),
	)

	router.Handle("GET /comments-count/{id}",
		endpoints.GetCommentCount(options.logger, options.commentService))
}

func addFeedRoutes(router *http.ServeMux, options options) {
	router.HandleFunc("GET /feed", endpoints.GetFeed(options.feedService))
}

func addAccountRoutes(router *http.ServeMux, options options) {
	router.Handle("POST /login",
		endpoints.Login(options.logger, options.accountService))

	router.Handle("GET /logout", endpoints.Logout())
}

func addAnnounceRoutes(router *http.ServeMux, options options) {
	router.Handle("GET /announce",
		endpoints.GetAnnounce(options.logger, options.announceService))

	router.Handle("POST /announce",
		middleware.Admin(
			endpoints.CreateAnnounce(options.logger, options.announceService),
		),
	)
	router.Handle("DELETE /announce",
		middleware.Admin(
			endpoints.DeleteAnnounce(options.logger, options.announceService),
		),
	)
}

func addProfileRoutes(router *http.ServeMux, options options) {
	router.Handle("GET /profile",
		endpoints.GetProfile(options.logger, options.profileService))
}

func addHealthRoutes(router *http.ServeMux, options options) {
	router.HandleFunc("GET /health", endpoints.Healthcheck(options.healthService))
}

func addViewRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "index", "", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	router.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if err := utils.RenderView(w, "post", "blog", idStr); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	router.Handle("GET /admin", middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "admin", "admin", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})))

	router.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if _, err := session.GetSession(r); err == nil {
			http.Redirect(w, r, "/admin", http.StatusMovedPermanently)
			return
		}
		if err := utils.RenderView(w, "login", "login", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	router.HandleFunc("GET /blog", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "blog", "blog", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
}
