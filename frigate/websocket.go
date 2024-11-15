package frigate

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
	"github.com/TandyTuscon/frigate-telegram-ws/common" // Import common for log and other utilities
)

// ListenWebSocket connects to the Frigate WebSocket and sends new events to the event channel
func ListenWebSocket(conf *config.Config, eventChan chan EventStruct) {
	for {
		// Establish WebSocket connection
		conn, _, err := websocket.DefaultDialer.Dial(conf.FrigateWebSocketURL, nil)
		if err != nil {
			common.LogError.Printf("Failed to connect to WebSocket: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		common.LogInfo.Println("Connected to Frigate WebSocket")

		// Listen for messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				common.LogError.Printf("WebSocket read error: %v. Reconnecting...", err)
				conn.Close()
				break // Exit the inner loop to reconnect
			}

			// Parse the WebSocket message
			var payload EventPayload
			if err := json.Unmarshal(message, &payload); err != nil {
				common.LogError.Printf("Failed to parse WebSocket message: %v. Raw message: %s", err, string(message))
				continue
			}

			// Handle "new" event types
			if payload.Type == "new" {
				common.LogInfo.Printf("Received new event: %s from camera: %s", payload.After.ID, payload.After.Camera)
				// Send the event to the event channel
				eventChan <- payload.After
			}
		}

		// Reconnect after connection is closed
		common.LogInfo.Println("WebSocket connection closed. Reconnecting...")
		time.Sleep(5 * time.Second)
	}
}
