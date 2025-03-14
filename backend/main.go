package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(hub *Hub, gm *GameManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	// Set the hub field when creating the client.
	client := &Client{
		conn: conn,
		// fortnite
		send: make(chan []byte),
		hub:  hub,
		gm:   gm,
	}

	hub.register <- client

	go client.writeMessages()
	client.readMessages()
}

func main() {
	hub := NewHub()
	go hub.Run()

	gm := NewGameManager(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, gm, w, r)
	})

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("ListenAndServe error: ", err)
	}
}
