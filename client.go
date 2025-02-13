package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	playerID string
	tableID  string
}

func (c *Client) readMessages(gm *GameManager) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read err: ", err)
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
		case "check":
			gm.HandleCheck(c)
		case "fold":
			gm.HandleFold(c)
		default:
			log.Println("unknown command type: ", gameMsg.Type)
		}
	}
}
