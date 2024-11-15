package main

import (
	"log"

	"frigate-telegram-ws/config"
	"frigate-telegram-ws/frigate"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	conf := config.LoadConfig("config.yml")

	bot, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		log.Printf("Error: %v", err)
		log.Printf("Failed to initialize bot: %v", err)
	}

	eventChan := make(chan frigate.EventStruct, 100)
	go frigate.ListenWebSocket(conf, eventChan)

	workerPool := frigate.NewWorkerPool(10)
	go workerPool.Start(bot, conf)

	for event := range eventChan {
		workerPool.AddTask(event)
	}
}
