package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	staticHandler "github.com/thiago-ssilva/zap/internal/api/handler/static"
	websocketHandler "github.com/thiago-ssilva/zap/internal/api/handler/websocket"
)

func SetupRouter(staticH *staticHandler.StaticHandler, websocketH *websocketHandler.WebsocketHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/static/*", staticH.StaticFiles)
	r.Get("/", staticH.Index)
	r.Get("/ws", websocketH.JoinRoom)

	return r
}
