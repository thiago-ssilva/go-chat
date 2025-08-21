package main

import (
	"log"
	"net/http"

	"github.com/thiago-ssilva/zap/internal/db"
	"github.com/thiago-ssilva/zap/internal/db/migrations"
	"github.com/thiago-ssilva/zap/internal/handler"
	"github.com/thiago-ssilva/zap/internal/repository"
	"github.com/thiago-ssilva/zap/internal/service"
	"github.com/thiago-ssilva/zap/internal/ws"
	"github.com/thiago-ssilva/zap/router"
)

func main() {
	//Database
	dbConn, err := db.NewDatabase()

	if err != nil {
		log.Fatalf("Could not initialize DB connection: %s", err)
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	//Run migrations
	if err := migrations.RunMigrations(dbConn); err != nil {
		log.Fatal(err)
	}

	// Set up repositories
	messagesRepo := repository.NewMessagesRepository(dbConn)

	// Set up Ws
	wsHub := ws.NewHub(messagesRepo)

	// Set up Services
	userService := service.NewUserService(wsHub)

	// Set up Handlers
	staticH := handler.NewStaticHandler()
	websocketH := handler.NewWebsocketHandler(wsHub, userService)
	usersH := handler.NewUserHandler(userService)

	go wsHub.Run()

	// Set up Server
	router := router.SetupRouter(staticH, websocketH, usersH)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Server starting on port %v\n", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
