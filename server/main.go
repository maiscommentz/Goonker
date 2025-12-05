package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Goonker/common"
	"Goonker/server/hub"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Server configuration constants
const (
	// Network configuration
	ServerPort       = ":8080"
	WsRoute          = "/ws"
	HandshakeTimeout = 5 * time.Second

	// Closure Reasons
	ErrExpectedJoin    = "Expected Join Packet"
	ErrFirstMustBeJoin = "First message must be 'join'"
	ErrInvalidPayload  = "Invalid Payload"
	ErrRoomIDRequired  = "Room ID required"
	ErrRoomFull        = "Room is full"
)

// main is the entry point of the server application.
func main() {
	// Register the WebSocket handler
	http.HandleFunc(WsRoute, wsHandler)

	// Start the server
	log.Printf("Starting server on port %s...", ServerPort)
	if err := http.ListenAndServe(ServerPort, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// wsHandler handles the initial HTTP upgrade and the application-layer handshake.
// Once the player is validated, control is passed to the Hub/Room.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket (skip verify for local dev)
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("Error upgrading websocket: %v", err)
		return
	}

	// Enforce timeout for handshake phase
	ctx, cancel := context.WithTimeout(r.Context(), HandshakeTimeout)
	defer cancel()

	// Read the first packet sent by the client
	var packet common.Packet
	if err := wsjson.Read(ctx, c, &packet); err != nil {
		log.Printf("Connection closed before join: %v", err)
		c.Close(websocket.StatusPolicyViolation, ErrExpectedJoin)
		return
	}

	// Validate that the first packet is a Join request
	if packet.Type != common.MsgJoin {
		log.Printf("First packet was not join: %s", packet.Type)
		c.Close(websocket.StatusPolicyViolation, ErrFirstMustBeJoin)
		return
	}

	// Parse the Join payload
	var joinData common.JoinPayload
	if err := json.Unmarshal(packet.Data, &joinData); err != nil {
		log.Printf("Invalid join payload: %v", err)
		c.Close(websocket.StatusProtocolError, ErrInvalidPayload)
		return
	}

	// Validate RoomID presence
	if joinData.RoomID == "" {
		c.Close(websocket.StatusPolicyViolation, ErrRoomIDRequired)
		return
	}

	// Let the Hub assign the player to a new or existing room
	room := hub.GlobalHub.CreateRoom(joinData.RoomID, joinData.IsBot)
	log.Printf("Client joining room '%s' (Bot: %v)", joinData.RoomID, joinData.IsBot)
	pid := room.AddPlayer(c)

	// Validation of assigned PlayerID, otherwise room is full
	if pid == common.Empty {
		log.Println("Room is full, rejecting client")
		c.Close(websocket.StatusPolicyViolation, ErrRoomFull)
	} else {
		log.Printf("Player assigned ID: %d in room %s", pid, joinData.RoomID)
	}
}
