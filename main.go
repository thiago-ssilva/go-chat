package main

import (
	"log"
	"net/http"

	staticHandler "github.com/thiago-ssilva/zap/internal/api/handler/static"
	"github.com/thiago-ssilva/zap/router"
)

func main() {

	staticH := staticHandler.NewStaticHandler()
	router := router.SetupRouter(staticH)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Server starting on port %v\n", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
