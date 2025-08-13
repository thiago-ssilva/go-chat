package main

import (
	"log"
	"net/http"

	"github.com/thiago-ssilva/todo-webapp/router"
)

func main() {
	router := router.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Server starting on port %v\n", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
