package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thiago-ssilva/zap/internal/service"
	"github.com/thiago-ssilva/zap/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketHandler struct {
	hub         *ws.Hub
	userService *service.UserService
}

func NewWebsocketHandler(hub *ws.Hub, userService *service.UserService) *WebsocketHandler {
	return &WebsocketHandler{
		hub:         hub,
		userService: userService,
	}
}

func (h *WebsocketHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	username := q.Get("username")

	if err := h.userService.ValidateUsername(username); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &ws.Client{
		Hub:      h.hub,
		Conn:     conn,
		Username: username,
		Send:     make(chan *ws.Message, 10),
	}

	h.hub.Register <- client

	go client.ReadMessage()
	go client.WriteMessage()
}
