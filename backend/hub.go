package main

import (
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

func (h *Hub) safeSend(ch chan []byte, msg []byte) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	ch <- msg
	return true
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
			// Only attempt to delete and close if still registered.
			if _, exists := h.clients[client]; exists {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			if message.Type == "public" {
				for client := range h.clients {
					if client.tableID == message.tableID {
						if !h.safeSend(client.send, message.content) {
							// If sending fails (likely due to a closed channel), remove the client.
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			} else if message.Type == "private" {
				for client := range h.clients {
					if client.playerID == message.playerID && client.tableID == message.tableID {
						if !h.safeSend(client.send, message.content) {
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
