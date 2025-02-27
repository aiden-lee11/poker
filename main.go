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
		log.Println("upgrade error: ", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte)}
	hub.register <- client

	// need to send handle join message?
	// dont think so i think for now on join were just creating a client
	// then the user will be prompted by frontend ui for the table they want to join and stack size they want which will read in via the client channel
	go client.writeMessages()
	client.readMessages(gm)
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
