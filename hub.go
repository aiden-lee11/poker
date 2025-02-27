package main

import (
	"fmt"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client)
			close(client.send)
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()

			if message.Type == "public" {
				for client := range h.clients {
					if client.tableID == message.tableID {
						select {
						case client.send <- message.content:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			} else if message.Type == "private" {
				for client := range h.clients {
					if client.playerID == message.playerID && client.tableID == message.tableID {
						fmt.Printf("client.playerID: %v\n", client.playerID)
						fmt.Printf("message.content: %v\n", message.content)
						select {
						case client.send <- message.content:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}

			h.mu.Unlock()
		}
	}
}
