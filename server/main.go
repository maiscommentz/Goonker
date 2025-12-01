package main

import (
	"log"
	"net/http"

	"Goonker/server/hub" // Adjust path to your module

	"nhooyr.io/websocket"
)

func main() {
	// 1. Create the Global Room Manager
	gameHub := hub.NewHub()

	// 2. Serve the static WASM files (from the /web folder)
	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	// 3. WebSocket Endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Accept the WebSocket connection (disable CORS for local dev)
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, 
		})
		if err != nil {
			log.Printf("Error accepting websocket: %v", err)
			return
		}

		// READ THE ROOM ID FROM URL: ws://localhost:8080/ws?room=ROOM_ID
		roomID := r.URL.Query().Get("room")
		if roomID == "" {
			roomID = "default" // Fallback if no code provided
		}

		log.Printf("New connection for room: %s", roomID)

		// Hand off connection to the Hub
		gameHub.HandleJoin(roomID, c, r.Context())
	})

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}