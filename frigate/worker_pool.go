package frigate

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
)

type WorkerPool struct {
	taskChan chan EventStruct
	wg       sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		taskChan: make(chan EventStruct, 100),
	}
}

func (wp *WorkerPool) Start(bot *tgbotapi.BotAPI, conf *config.Config) {
	for i := 0; i < 10; i++ { // Adjust worker count
		wp.wg.Add(1)
		go wp.worker(bot, conf)
	}
}

func (wp *WorkerPool) worker(bot *tgbotapi.BotAPI, conf *config.Config) {
	defer wp.wg.Done()
	for event := range wp.taskChan {
		ProcessEvent(event, bot, conf)
	}
}

func (wp *WorkerPool) AddTask(event EventStruct) {
	wp.taskChan <- event
}

func (wp *WorkerPool) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()
}
