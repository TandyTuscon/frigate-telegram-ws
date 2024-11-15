package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug              bool                     `yaml:"debug"`
	TelegramBotToken   string                   `yaml:"telegram_bot_token"`
	TelegramChatID     int64                    `yaml:"telegram_chat_id"`
	FrigateURL         string                   `yaml:"frigate_url"`
	FrigateExternalURL string                   `yaml:"frigate_external_url"`
	FrigateWebSocketURL string                  `yaml:"frigate_websocket_url"`
	Cameras            map[string]CameraConfig `yaml:"cameras"`
	MessageTemplates   MessageTemplates        `yaml:"message_templates"`
	SnapshotPath       string                  `yaml:"snapshot_path"`
	ClipPath           string                  `yaml:"clip_path"`
	TimeZone           string                  `yaml:"time_zone"`
}

// CameraConfig defines the configuration for a camera
type CameraConfig struct {
	Label []string               `yaml:"label"`
	Zone  map[string]ZoneConfig  `yaml:"zone"`
	Score ScoreConfig            `yaml:"score"`
	Length LengthConfig          `yaml:"length"`
	Severity SeverityConfig      `yaml:"severity"`
	Media MediaConfig            `yaml:"media"`
	Time TimeConfig              `yaml:"time"`
}

// ZoneConfig defines nested conditions for specific zones
type ZoneConfig struct {
	Label map[string]LabelConfig `yaml:"label"`
	Severity SeverityConfig      `yaml:"severity"`
}

// LabelConfig defines conditions for specific labels within a zone
type LabelConfig struct {
	Sublabel []string  `yaml:"sublabel"`
	Score    ScoreConfig `yaml:"score"`
}

// ScoreConfig defines the minimum and maximum score thresholds
type ScoreConfig struct {
	MinScore float64 `yaml:"min_score"`
	MaxScore float64 `yaml:"max_score,omitempty"`
}

// LengthConfig defines minimum and maximum length thresholds
type LengthConfig struct {
	MinLength float64 `yaml:"min_length"`
	MaxLength float64 `yaml:"max_length"`
}

// SeverityConfig defines severity-related filters
type SeverityConfig struct {
	Reviewed bool   `yaml:"reviewed"`
	Severity string `yaml:"severity"`
}

// MediaConfig defines media-related filters
type MediaConfig struct {
	HasSnapshot      bool `yaml:"has_snapshot"`
	HasClip          bool `yaml:"has_clip"`
	IncludeThumbnails bool `yaml:"include_thumbnails"`
}

// TimeConfig defines time-related filters
type TimeConfig struct {
	Before    int64  `yaml:"before"`
	After     int64  `yaml:"after"`
	Timezone  string `yaml:"timezone"`
	TimeRange string `yaml:"time_range"`
}

// MessageTemplates defines templates for notifications
type MessageTemplates struct {
	TitleTemplate string            `yaml:"title_template"`
	BodyTemplate  string            `yaml:"body_template"`
	Fields        []string          `yaml:"fields"`
	DisplayNames  DisplayNames      `yaml:"display_names"`
	CustomMessages []CustomMessage  `yaml:"custom_messages"`
}

// DisplayNames defines user-friendly names for cameras, zones, labels, and sublabels
type DisplayNames struct {
	Cameras map[string]string `yaml:"cameras"`
	Zones   map[string]string `yaml:"zones"`
	Labels  map[string]string `yaml:"labels"`
	Sublabels map[string]string `yaml:"sublabels"`
}

// CustomMessage defines condition-based message customization
type CustomMessage struct {
	Camera     string            `yaml:"camera,omitempty"`
	Zone       string            `yaml:"zone,omitempty"`
	Label      string            `yaml:"label,omitempty"`
	Sublabel   string            `yaml:"sublabel,omitempty"`
	Text       string            `yaml:"text"`
	Formatting map[string]FormatOptions `yaml:"formatting"`
}

// FormatOptions defines formatting options for fields
type FormatOptions struct {
	CapitalizeFirstLetter bool   `yaml:"capitalize_first_letter"`
	CapitalizeAllWords    bool   `yaml:"capitalize_all_words"`
	AllCaps               bool   `yaml:"all_caps"`
	Bold                  bool   `yaml:"bold"`
	Italics               bool   `yaml:"italics"`
	Emoji                 string `yaml:"emoji,omitempty"`
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
