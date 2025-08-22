package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thiago-ssilva/go-chat/internal/handler"
)

func SetupRouter(staticH *handler.StaticHandler, websocketH *handler.WebsocketHandler, userH *handler.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/static/*", staticH.StaticFiles)
	r.Get("/", staticH.Index)
	r.Get("/ws", websocketH.JoinRoom)
	r.Get("/api/users/validate/username", userH.ValidateUsername)

	return r
}
