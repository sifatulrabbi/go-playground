package tgd

import (
	"fmt"
)

type ChatHub struct {
	// Registered clients.
	clients map[string]*Client
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

func (h *ChatHub) start() {
	fmt.Println("Starting the chat hub...")
	for {
		select {
		case client := <-h.register:
			if client != nil {
				fmt.Println("New user connected:", client.id)
				h.clients[client.id] = client
			}

		case client := <-h.unregister:
			if client != nil {
				fmt.Println("User disconnected:", client.id)
				delete(h.clients, client.id)
			}

		case msg := <-h.broadcast:
			for _, client := range h.clients {
				client.send <- msg
			}
		}
	}
}

func NewChatHub() *ChatHub {
	cb := ChatHub{
		clients:    map[string]*Client{},
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	return &cb
}
