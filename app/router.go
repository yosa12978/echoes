package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/handlers"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

func NewRouter(ctx context.Context) http.Handler {
	postRepo := repos.NewPostPostgres()
	linkRepo := repos.NewLinkPostgres()
	announceRepo := repos.NewAnnounce()
	profileRepo, err := repos.NewProfileJson("./static/profile.json")
	if err != nil {
		log.Println(err.Error())
	}
	//postRepo.Seed(ctx)
	//announceRepo.Create("*beep* *boop* new announce *beep* *boop*")

	postHandler := handlers.NewPost(postRepo)
	linkHandler := handlers.NewLink(linkRepo)
	announceHandler := handlers.NewAnnounce(announceRepo)
	profileHandler := handlers.NewProfile(profileRepo)

	router := mux.NewRouter()
	router.StrictSlash(true)

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

	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.RenderTemplate(w, "admin", nil); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}).Methods("GET")

	hateoas := router.PathPrefix("/hateoas").Subrouter()
	hateoas.Handle("/posts", postHandler.GetPosts(ctx)).Methods("GET")
	hateoas.Handle("/posts/{id}", postHandler.GetPostById(ctx)).Methods("GET")
	hateoas.Handle("/posts", postHandler.CreatePost(ctx)).Methods("POST")
	hateoas.Handle("/posts", postHandler.DeletePost(ctx)).Methods("DELETE")
	hateoas.Handle("/post-pin", postHandler.PinPost(ctx)).Methods("PATCH")

	hateoas.Handle("/links", linkHandler.GetAll(ctx)).Methods("GET")
	hateoas.Handle("/links-admin", linkHandler.GetAdmin(ctx)).Methods("GET")
	hateoas.Handle("/links", linkHandler.Create(ctx)).Methods("POST")
	hateoas.Handle("/links/{id}", linkHandler.Delete(ctx)).Methods("DELETE")

	hateoas.Handle("/announce", announceHandler.Get(ctx)).Methods("GET")
	hateoas.Handle("/profile", profileHandler.Get(ctx)).Methods("GET")
	return router
}
