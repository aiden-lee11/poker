package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	playerID  string
	tableID   string
	closeOnce sync.Once // ensures Close() only runs once
}

// Close safely unregisters the client and closes the connection.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		// Unregister client if the hub and its channel are non-nil.
		if c.hub != nil && c.hub.unregister != nil {
			// If unregister might block, you can use a non-blocking send:
			select {
			case c.hub.unregister <- c:
			default:
			}
		}
		// Close the send channel.
		close(c.send)
		// Finally, close the websocket connection.
		c.conn.Close()
	})
}

// readMessages reads messages from the client.
// It defers the Close() to ensure that cleanup happens only once.
func (c *Client) readMessages(gm *GameManager) {
	// Ensure cleanup happens when this function exits.
	defer c.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			break
		}

		var gameMsg GameMessage
		if err := json.Unmarshal(message, &gameMsg); err != nil {
			log.Println("json err:", err)
			continue
		}

		switch gameMsg.Type {
		case "join":
			gm.HandleJoin(c, gameMsg.Payload)
		case "bet":
			gm.HandleBet(c, gameMsg.Payload)
		case "call":
			gm.HandleCall(c)
		case "check":
			gm.HandleCheck(c)
		case "fold":
			gm.HandleFold(c)
		case "init":
			gm.HandleInitGame(c)
		default:
			log.Println("unknown command type:", gameMsg.Type)
		}
	}
}

// writeMessages writes outgoing messages to the client.
func (c *Client) writeMessages() {
	defer c.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
