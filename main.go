package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"poker/table"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	JoinGame      = "join"
	Bet           = "bet"
	Check         = "check"
	Fold          = "fold"
	DealCards     = "deal"
	EvaluateHands = "evaluate"
)

func initHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if client != message.sender {
					client.send <- message.content
				}

			}
			h.mu.Unlock()
		}

	}
}

func (c *Client) readMessages(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error: ", err)
			break
		}

		var gameMessage GameMessage
		err = json.Unmarshal(message, &gameMessage)

		if err != nil {
			log.Println("json read error: ", err)
			break
		}
		// JoinGame      = "join"
		// Bet           = "bet"
		// Check         = "check"
		// Fold          = "fold"
		// DealCards     = "deal"
		// EvaluateHands = "evaluate"
		var playerID string
		if payload, ok := gameMessage.Payload.(map[string]interface{}); ok {
			if id, ok := payload["playerID"].(string); ok {
				playerID = id
			} else {
				log.Println("no playerid in payload", gameMessage)
			}
		}

		switch gameMessage.Type {
		case JoinGame:
			// do something need client to join game func idk
			c.handleJoinGame(gameMessage.Payload)
		case DealCards:
			c.handleDealCards()
		case EvaluateHands:
			c.handleEvaluateHands()
		case Bet:
			c.handleBet(playerID, gameMessage.Payload)
		case Check:
			c.handleCheck(playerID, gameMessage.Payload)
		case Fold:
			c.handleFold(playerID, gameMessage.Payload)
		default:
			log.Println("how did you get here lol", gameMessage.Type)
		}
	}
}

func (c *Client) handleJoinGame(payload interface{}) {
	var stackSize int
	if payload, ok := payload.(map[string]interface{}); ok {
		if payloadStackSize, ok := payload["stackSize"].(int); ok {
			stackSize = payloadStackSize
		} else {
			log.Println("no playerid in payload", payload)
		}
	}

	playerID := fmt.Sprintf("player-%d", len(c.hub.clients)+1)
	playerData := table.Player{
		StackSize: stackSize,
		PlayerID:  playerID,
	}

	c.playerID = playerID
	c.playerData = &playerData

}

func (c *Client) handleBet(playerID string, payload interface{}) {

}

func (c *Client) writeMessages() {
	defer c.conn.Close()
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error: ", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte)}
	hub.register <- client

	go client.writeMessages()
	client.readMessages(hub)
}

func main() {
	hub := initHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("ListenAndServe error: ", err)
	}
}
