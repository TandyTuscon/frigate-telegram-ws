func ListenWebSocket(conf *config.Config, eventChan chan EventStruct) {
	conn, _, err := websocket.DefaultDialer.Dial(conf.FrigateWebSocketURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to WebSocket")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error: %v", err)
			time.Sleep(5 * time.Second)
			continue
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
