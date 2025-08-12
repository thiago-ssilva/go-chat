package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

var startTime time.Time

func main() {
	startTime = time.Now()

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.Handle("/health", new(healthHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
		Uptime string `json:"uptime"`
	}{
		Status: "healthy",
		Uptime: startTime.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
