package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/cache"
	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/handlers"
	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/middleware"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/session"
	"github.com/yosa12978/echoes/utils"
)

func NewRouter(ctx context.Context) http.Handler {
	logger := logging.New("app.NewRouter")

	postRepo := repos.NewPostPostgres()
	postService := services.NewPost(
		postRepo,
		cache.NewRedisCache(ctx),
		logging.New("postService"),
		repos.NewPostSearcherPostgres(),
	)

	linkRepo := repos.NewLinkPostgres()
	linkService := services.NewLink(
		linkRepo,
		cache.NewRedisCache(ctx),
		logging.New("linkService"),
	)

	commentRepo := repos.NewCommentPostgres()
	commentService := services.NewComment(
		commentRepo,
		postService,
		cache.NewRedisCache(ctx),
		logging.New("commentService"),
	)

	announceRepo := repos.NewAnnounceCache(cache.NewRedisCache(ctx))
	announceService := services.NewAnnounce(
		announceRepo,
		logging.New("announceService"),
	)

	accountRepo := repos.NewAccountPostgres()
	accountService := services.NewAccount(accountRepo)
	accountService.Seed(ctx)

	profileRepo, err := repos.NewProfileJson("./assets/profile.json")
	if err != nil {
		logger.Error(err)
	}
	profileService := services.NewProfile(profileRepo)

	feedService := services.NewFeedService(postService)

	postHandler := handlers.NewPost(postService, logging.New("postHandler"))
	linkHandler := handlers.NewLink(linkService, logging.New("linkHandler"))
	announceHandler := handlers.NewAnnounce(announceService, logging.New("announceHandler"))
	profileHandler := handlers.NewProfile(profileService, logging.New("profileHandler"))
	accountHandler := handlers.NewAccount(accountService, logging.New("accountHandler"))
	commentHandler := handlers.NewComment(commentService, logging.New("commentHandler"))
	feedHandler := handlers.NewFeedHandler(feedService)

	//latencyLogger := middleware.Logger(logging.New("request"))
	router := mux.NewRouter()
	router.StrictSlash(true)

	//router.Use(latencyLogger)

	RegisterBasicHandler(ctx, router)

	hateoas := router.PathPrefix("/api").Subrouter()

	RegisterLinkHandler(ctx, linkHandler, hateoas)
	RegisterAccountHandler(ctx, accountHandler, hateoas)
	RegisterPostHandler(ctx, postHandler, hateoas)
	RegisterProfileHandler(ctx, profileHandler, hateoas)
	RegisterAnnounceHandler(ctx, announceHandler, hateoas)
	RegisterCommentHandler(ctx, commentHandler, hateoas)
	RegisterFeedHandler(feedHandler, router)

	return router
}

func RegisterLinkHandler(ctx context.Context, handler handlers.Link, router *mux.Router) {
	router.Handle("/links", handler.GetAll(ctx)).Methods("GET")
	router.Handle("/links-admin", middleware.Admin(handler.GetAdmin(ctx))).Methods("GET")
	router.Handle("/links", middleware.Admin(handler.Create(ctx))).Methods("POST")
	router.Handle("/links/{id}", middleware.Admin(handler.Delete(ctx))).Methods("DELETE")

	router.Handle("/portal/{id}", handler.Portal(ctx)).Methods("GET")
}

func RegisterFeedHandler(handler handlers.Feed, router *mux.Router) {
	router.HandleFunc("/feed", handler.GetFeed()).Methods("GET")
}

func RegisterPostHandler(ctx context.Context, handler handlers.Post, router *mux.Router) {
	router.Handle("/posts", handler.GetPosts(ctx)).Methods("GET")
	router.Handle("/posts/{id}", handler.GetPostById(ctx)).Methods("GET")

	router.Handle("/posts", middleware.Admin(handler.CreatePost(ctx))).Methods("POST")
	router.Handle("/posts", middleware.Admin(handler.DeletePost(ctx))).Methods("DELETE")
	router.Handle("/post-pin", middleware.Admin(handler.PinPost(ctx))).Methods("PATCH")
}

func RegisterCommentHandler(ctx context.Context, handler handlers.Comment, router *mux.Router) {
	router.Handle("/comments", handler.GetPostComments(ctx)).Methods("GET")
	router.Handle("/comments", handler.CreateComment(ctx)).Methods("POST")
	router.Handle("/comments/{id}", middleware.Admin(handler.DeleteComment(ctx))).Methods("DELETE")
	router.Handle("/comments-count/{id}", handler.GetCommentCount(ctx)).Methods("GET")
}

func RegisterAnnounceHandler(ctx context.Context, handler handlers.Announce, router *mux.Router) {
	router.Handle("/announce", handler.Get(ctx)).Methods("GET")
	router.Handle("/announce", middleware.Admin(handler.Create(ctx))).Methods("POST")
	router.Handle("/announce", middleware.Admin(handler.Delete(ctx))).Methods("DELETE")
}

func RegisterAccountHandler(ctx context.Context, handler handlers.Account, router *mux.Router) {
	router.Handle("/login", handler.Login(ctx)).Methods("POST")
	router.Handle("/logout", handler.Logout(ctx)).Methods("GET")
}

func RegisterProfileHandler(ctx context.Context, handler handlers.Profile, router *mux.Router) {
	router.Handle("/profile", handler.Get(ctx)).Methods("GET")
}

func RegisterBasicHandler(ctx context.Context, router *mux.Router) {
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "index", "", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		if err := utils.RenderView(w, "post", "", idStr); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	// router.HandleFunc("/portal", func(w http.ResponseWriter, r *http.Request) {
	// 	url := r.URL.Query().Get("url")
	// 	if url == "" {
	// 		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	// 		return
	// 	}
	// 	http.Redirect(w, r, url, http.StatusMovedPermanently)
	// }).Methods("GET")

	// simplify this (cuz it looks terrible)
	router.Handle("/admin", middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "admin", "admin", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}))).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if _, err := session.GetInfo(r); err == nil {
			http.Redirect(w, r, "/admin", http.StatusMovedPermanently)
			return
		}
		if err := utils.RenderView(w, "login", "login", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	router.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderView(w, "blog", "blog", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	// Healthchecks
	healthService := services.NewHealthService(
		logging.New("healthchecks"),
		data.NewPgPinger(),
		data.NewRedisPinger(ctx),
	)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := healthService.Healthcheck(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	}).Methods("GET")
}
