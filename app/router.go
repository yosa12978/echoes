package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
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
	linkRepo := repos.NewLinkPostgres()
	announceRepo := repos.NewAnnounce()
	accountRepo := repos.NewAccountPostgres()
	postService := services.NewPost(repos.NewPostPostgres())
	//postService.Seed(ctx)
	commentService := services.NewComment(repos.NewCommentPostgres(), postService)
	profileRepo, err := repos.NewProfileJson("./static/profile.json")
	if err != nil {
		logger.Error(err)
	}

	postHandler := handlers.NewPost(postRepo)
	linkHandler := handlers.NewLink(linkRepo)
	announceHandler := handlers.NewAnnounce(announceRepo)
	profileHandler := handlers.NewProfile(profileRepo)
	accountHandler := handlers.NewAccount(accountRepo)
	commentHandler := handlers.NewComment(commentService)

	router := mux.NewRouter()
	router.StrictSlash(true)

	RegisterBasicHandler(ctx, router)

	hateoas := router.PathPrefix("/hateoas").Subrouter()

	RegisterLinkHandler(ctx, linkHandler, hateoas)
	RegisterAccountHandler(ctx, accountHandler, hateoas)
	RegisterPostHandler(ctx, postHandler, hateoas)
	RegisterProfileHandler(ctx, profileHandler, hateoas)
	RegisterAnnounceHandler(ctx, announceHandler, hateoas)
	RegisterCommentHandler(ctx, commentHandler, hateoas)

	return router
}

func RegisterLinkHandler(ctx context.Context, handler handlers.Link, router *mux.Router) {
	router.Handle("/links", handler.GetAll(ctx)).Methods("GET")
	router.Handle("/links-admin", middleware.Admin(handler.GetAdmin(ctx))).Methods("GET")
	router.Handle("/links", middleware.Admin(handler.Create(ctx))).Methods("POST")
	router.Handle("/links/{id}", middleware.Admin(handler.Delete(ctx))).Methods("DELETE")
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
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderTemplate(w, "index", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		if err := utils.RenderTemplate(w, "post", idStr); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	router.HandleFunc("/portal", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}).Methods("GET")

	// simplify this (cuz it looks terrible)
	router.Handle("/admin", middleware.Admin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderTemplate(w, "admin", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}))).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if _, err := session.GetInfo(r); err == nil {
			http.Redirect(w, r, "/admin", http.StatusMovedPermanently)
			return
		}
		if err := utils.RenderTemplate(w, "login", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")
}
