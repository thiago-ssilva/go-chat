package ws

import (
	"context"
	"log"

	"github.com/thiago-ssilva/zap/internal/repository"
)

type Hub struct {
	clients      map[*Client]bool
	Register     chan *Client
	Unregister   chan *Client
	Broadcast    chan *Message
	messagesRepo *repository.MessagesRepository
}

func NewHub(messagesRepo *repository.MessagesRepository) *Hub {
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
			//goroutine to get current saved messages
			go func() {
				messages, err := h.messagesRepo.GetAllMessages(context.Background())

				if err != nil {
					log.Printf("Failed to load messages: %v", err)
				}

				for _, msg := range messages {
					wsMsg := &Message{
						Content:  msg.Content,
						Username: msg.Username,
					}

					client.Send <- wsMsg
				}
			}()

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			// goroutine to persist message
			go func(msg *Message) {
				messageDb := &repository.Message{
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

func (h *Hub) IsUsernameTaken(username string) bool {
	for client := range h.clients {
		if username == client.Username {
			return true
		}
	}

	return false
}
