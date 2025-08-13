package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.HandleFunc("/static/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fileServer.ServeHTTP(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	return r
}
