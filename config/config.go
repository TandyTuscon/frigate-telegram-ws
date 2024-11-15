package config

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug                bool               `yaml:"debug"`
	TelegramBotToken     string             `yaml:"telegram_bot_token"`
	TelegramChatID       int64              `yaml:"telegram_chat_id"`
	FrigateURL           string             `yaml:"frigate_url"`
	FrigateWebSocketURL  string             `yaml:"frigate_websocket_url"` // Add WebSocket URL
	FrigateExternalURL   string             `yaml:"frigate_external_url"`  // External URL for media sharing
	Cameras              map[string]CameraConfig `yaml:"cameras"`
	MessageTemplates     MessageTemplates   `yaml:"message"`
	WorkerPoolSize       int                `yaml:"worker_pool_size"`       // Worker pool size configuration
	TaskQueueBufferSize  int                `yaml:"task_queue_buffer_size"` // Task queue buffer size
	TimeFormat           string             `yaml:"time_format"`            // Custom time format
}

// CameraConfig defines per-camera configurations
type CameraConfig struct {
	Label  []string `yaml:"label"`
	Zone   []string `yaml:"zone"`
	Score  struct {
		MinScore float64 `yaml:"min_score"`
		MaxScore float64 `yaml:"max_score"`
	} `yaml:"score"`
	Length struct {
		MinLength float64 `yaml:"min_length"`
		MaxLength float64 `yaml:"max_length"`
	} `yaml:"length"`
}

// MessageTemplates define templates for event messages
type MessageTemplates struct {
	TitleTemplate string `yaml:"title_template"`
	BodyTemplate  string `yaml:"body_template"`
}

// LoadConfig reads the YAML configuration file and populates the Config struct
func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	var conf Config
	if err := yaml.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Validate required fields
	if conf.TelegramBotToken == "" {
		log.Fatalf("TelegramBotToken is required in config.yml")
	}
	if conf.FrigateURL == "" {
		log.Fatalf("FrigateURL is required in config.yml")
	}
	if conf.TelegramChatID == 0 {
		log.Fatalf("TelegramChatID is required in config.yml")
	}

	// Apply default values for optional fields
	if conf.WorkerPoolSize == 0 {
		conf.WorkerPoolSize = 10 // Default worker pool size
	}
	if conf.TaskQueueBufferSize == 0 {
		conf.TaskQueueBufferSize = 100 // Default buffer size
	}
	if conf.TimeFormat == "" {
		conf.TimeFormat = "standard" // Default time format
	}

	return &conf
}
