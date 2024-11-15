package main

import (
	"log"

	"github.com/TandyTuscon/frigate-telegram/config"
	"github.com/TandyTuscon/frigate-telegram/frigate"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load configuration
	conf := config.LoadConfig("config.yml")

	// Initialize the Telegram bot
	bot, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Create a channel to handle events
	eventChan := make(chan frigate.EventStruct, 100)

	// Start listening to Frigate WebSocket
	go frigate.ListenWebSocket(conf, eventChan)

	// Initialize and start the worker pool
	workerPool := frigate.NewWorkerPool(10)
	go workerPool.Start(bot, conf)

	// Process events from the channel
	for event := range eventChan {
		workerPool.AddTask(event)
	}
}
