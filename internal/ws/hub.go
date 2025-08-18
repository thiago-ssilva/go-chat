package ws

import (
	"context"
	"log"

	"github.com/thiago-ssilva/zap/internal/repositories"
)

type Hub struct {
	clients      map[*Client]bool
	Register     chan *Client
	Unregister   chan *Client
	Broadcast    chan *Message
	messagesRepo *repositories.MessagesRepository
}

func NewHub(messagesRepo *repositories.MessagesRepository) *Hub {
	return &Hub{
		Broadcast:    make(chan *Message, 5),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		messagesRepo: messagesRepo,
	}
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			go func(msg *Message) {
				messageDb := &repositories.Message{
					Content:  message.Content,
					Username: message.Username,
				}
				if _, err := h.messagesRepo.CreateMessage(context.Background(), messageDb); err != nil {
					log.Printf("Failed to persist message: %v", err)
				}
			}(message)

			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}

		}
	}
}
