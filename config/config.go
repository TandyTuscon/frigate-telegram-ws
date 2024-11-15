package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug              bool
	TelegramBotToken   string
	TelegramChatID     int64
	FrigateURL         string
	FrigateExternalURL string `yaml:"frigate_external_url"` // Added for external Frigate URL
	Cameras            map[string]CameraConfig
	MessageTemplates   MessageTemplates
	SnapshotPath       string `yaml:"snapshot_path"`       // Added for snapshot saving path
	ClipPath           string `yaml:"clip_path"`          // Added for clip saving path
	TimeZone           string `yaml:"time_zone"`          // Added for event timezone
}

type CameraConfig struct {
	Label []string `yaml:"label"`
	Score struct {
		MinScore float64 `yaml:"min_score"`
	} `yaml:"score"`
}

type MessageTemplates struct {
	TitleTemplate string `yaml:"title_template"`
	BodyTemplate  string `yaml:"body_template"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error: Failed to open config file: %v", err)
		return nil
	}
	defer file.Close()

	var conf Config
	if err := yaml.NewDecoder(file).Decode(&conf); err != nil {
		log.Printf("Error: Failed to parse config: %v", err)
		return nil
	}

	// Log loaded configuration if debug mode is enabled
	if conf.Debug {
		log.Printf("Configuration loaded successfully: %+v", conf)
	}

	return &conf
}
