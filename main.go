package main

import (
	"log"
	"net/http"

	staticHandler "github.com/thiago-ssilva/zap/internal/api/handler/static"
	websocketHandler "github.com/thiago-ssilva/zap/internal/api/handler/websocket"
	"github.com/thiago-ssilva/zap/internal/db"
	"github.com/thiago-ssilva/zap/internal/db/migrations"
	"github.com/thiago-ssilva/zap/internal/repositories"
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
	messagesRepo := repositories.NewMessagesRepository(dbConn)

	// Set up Ws
	wsHub := ws.NewHub(messagesRepo)

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
