package main

import (
	"log"

	"github.com/TandyTuscon/frigate-telegram-ws/config"  // Updated import path
	"github.com/TandyTuscon/frigate-telegram-ws/frigate" // Updated import path

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load configuration from config.yml
	conf := config.LoadConfig("config.yml")

	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Channel for processing events
	eventChan := make(chan frigate.EventStruct, 100)

	// Start listening to the WebSocket for Frigate events
	go frigate.ListenWebSocket(conf, eventChan)

	// Initialize and start the worker pool
	workerPool := frigate.NewWorkerPool(10)
	go workerPool.Start(bot, conf)

	// Dispatch events to the worker pool
	for event := range eventChan {
		workerPool.AddTask(event)
	}
}
