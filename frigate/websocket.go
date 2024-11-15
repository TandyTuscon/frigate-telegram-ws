package frigate

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
)

// EventPayload defines the structure of WebSocket payloads
type EventPayload struct {
	Type  string      `json:"type"`
	After EventStruct `json:"after"`
}

// ListenWebSocket connects to the Frigate WebSocket and processes events
func ListenWebSocket(conf *config.Config, eventChan chan EventStruct) {
	conn, _, err := websocket.DefaultDialer.Dial(conf.FrigateWebSocketURL, nil)
	if err != nil {
		log.Printf("Error: Failed to connect to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Connected to WebSocket")

	reconnectAttempts := 0
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error: %v. Attempting to reconnect...", err)
			time.Sleep(5 * time.Second)
			reconnectAttempts++
			if reconnectAttempts > 10 {
				log.Printf("Exceeded maximum reconnect attempts. Exiting WebSocket listener.")
				return
			}
			conn, _, err = websocket.DefaultDialer.Dial(conf.FrigateWebSocketURL, nil)
			if err != nil {
				log.Printf("Reconnect failed: %v", err)
				continue
			}
			reconnectAttempts = 0
			log.Println("Reconnected to WebSocket.")
		}

		var payload EventPayload
		if err := json.Unmarshal(message, &payload); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if payload.Type == "new" {
			eventChan <- payload.After
		}
	}
}
