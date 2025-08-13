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

func (h *WebsocketHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := &ws.Client{
		Hub:  h.hub,
		Conn: conn,
	}

}
