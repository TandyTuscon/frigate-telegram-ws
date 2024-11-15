package frigate

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
)

// WorkerPool represents a pool of workers that process events
type WorkerPool struct {
	taskChan chan EventStruct
	wg       sync.WaitGroup
	stopChan chan struct{} // Channel to signal workers to stop
}

// NewWorkerPool initializes a new worker pool with the specified size
func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		taskChan: make(chan EventStruct, 100), // Task queue with a buffer
		stopChan: make(chan struct{}),        // Signal channel for stopping workers
	}
}

// Start initializes the workers in the pool
func (wp *WorkerPool) Start(workerCount int, bot *tgbotapi.BotAPI, conf *config.Config) {
	log.Printf("Starting %d workers", workerCount)
	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(bot, conf, i)
	}
}

// worker processes tasks from the task channel
func (wp *WorkerPool) worker(bot *tgbotapi.BotAPI, conf *config.Config, id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case event := <-wp.taskChan:
			// Processing the event
			log.Printf("Worker %d processing event: %s", id, event.ID)
			ProcessEvent(event, bot, conf)

		case <-wp.stopChan:
			// Graceful shutdown signal received
			log.Printf("Worker %d stopping", id)
			return
		}
	}
}

// AddTask adds a new task to the task channel
func (wp *WorkerPool) AddTask(event EventStruct) {
	select {
	case wp.taskChan <- event:
		log.Printf("Task added to the queue: %s", event.ID)
	default:
		// If the task queue is full, it logs the event and drops it
		log.Printf("Task queue is full, dropping event: %s", event.ID)
	}
}

// Stop signals all workers to stop and waits for them to finish
func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool")
	close(wp.stopChan) // Signal workers to stop
	wp.wg.Wait()       // Wait for all workers to finish
	log.Println("Worker pool stopped")
}
