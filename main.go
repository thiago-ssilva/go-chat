package main

import (
	"log"
	"net/http"

	staticHandler "github.com/thiago-ssilva/zap/internal/api/handler/static"
	websocketHandler "github.com/thiago-ssilva/zap/internal/api/handler/websocket"
	"github.com/thiago-ssilva/zap/internal/ws"
	"github.com/thiago-ssilva/zap/router"
)

func main() {
	// Set up Ws
	wsHub := ws.NewHub()

	go wsHub.Run()

	// Set up Handlers
	staticH := staticHandler.NewStaticHandler()
	websocketH := websocketHandler.NewWebsocketHandler(wsHub)

	// Set up Server
	router := router.SetupRouter(staticH, websocketH)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Server starting on port %v\n", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
