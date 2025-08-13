package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thiago-ssilva/zap/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketHandler struct {
	hub *ws.Hub
}

func NewWebsocketHandler(hub *ws.Hub) *WebsocketHandler {
	return &WebsocketHandler{
		hub: hub,
	}
}

func (h *WebsocketHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	q := r.URL.Query()
	username := q.Get("username")

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
