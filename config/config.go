package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug            bool
	TelegramBotToken string
	TelegramChatID   int64
	FrigateURL       string
	Cameras          map[string]CameraConfig
	MessageTemplates MessageTemplates
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
