package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yosa12978/echoes/handlers"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/utils"
)

func NewRouter(ctx context.Context) http.Handler {
	postRepo := repos.NewPostPostgres()
	//linkRepo := repos.NewLinkPostgres()

	postHandler := handlers.NewPost(postRepo)

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RenderTemplate(w, "index", nil)
	}).Methods("GET")

	router.Handle("/posts", postHandler.GetPosts(ctx)).Methods("GET")
	router.Handle("/posts/{id}", postHandler.GetPostById(ctx)).Methods("GET")
	router.Handle("/posts", postHandler.CreatePost(ctx)).Methods("POST")
	router.Handle("/posts/{id}", postHandler.DeletePost(ctx)).Methods("DELETE")

	return router
}
