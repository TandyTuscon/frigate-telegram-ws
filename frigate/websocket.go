package frigate

import (
    "encoding/json"
    "time"

    "github.com/gorilla/websocket"
    "github.com/TandyTuscon/frigate-telegram-ws/config"
    "github.com/TandyTuscon/frigate-telegram-ws/frigate"
)

// ListenWebSocket connects to the Frigate WebSocket and sends new events to the event channel
func ListenWebSocket(conf *config.Config, eventChan chan EventStruct) {
    for {
        // Establish WebSocket connection
        conn, _, err := websocket.DefaultDialer.Dial(conf.FrigateWebSocketURL, nil)
        if err != nil {
            log.Error.Printf("Failed to connect to WebSocket: %v. Retrying in 5 seconds...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        log.Info.Println("Connected to Frigate WebSocket")

        // Listen for messages
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                log.Error.Printf("WebSocket read error: %v. Reconnecting...", err)
                conn.Close()
                break // Exit the inner loop to reconnect
            }

            // Parse the WebSocket message
            var payload EventPayload
            if err := json.Unmarshal(message, &payload); err != nil {
                log.Error.Printf("Failed to parse WebSocket message: %v. Raw message: %s", err, string(message))
                continue
            }

            // Handle "new" event types
            if payload.Type == "new" {
                log.Info.Printf("Received new event: %s from camera: %s", payload.After.ID, payload.After.Camera)
                eventChan <- payload.After
            }
        }

        // Reconnect after connection is closed
        log.Info.Println("WebSocket connection closed. Reconnecting...")
        time.Sleep(5 * time.Second)
    }
}
