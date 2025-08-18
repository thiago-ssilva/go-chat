package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	//Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	//Max message size allowed from peer.
	maxMessageSize = 512

	//Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	//Time to send pings to peer.
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan *Message
	Username string
}

func (c *Client) ReadMessage() {

	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, m, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			Username: c.Username,
		}

		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			n := len(c.Send)

			batch := struct {
				Messages []*Message `json:"messages"`
			}{
				Messages: make([]*Message, 0, n+1),
			}

			batch.Messages = append(batch.Messages, message)

			for range n {
				queuedMsg := <-c.Send

				batch.Messages = append(batch.Messages, queuedMsg)
			}

			if err := c.Conn.WriteJSON(batch); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
