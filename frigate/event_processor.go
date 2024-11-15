package frigate

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
	"github.com/TandyTuscon/frigate-telegram-ws/common" // Import common package
)

// ProcessEvent processes an event and sends a message to Telegram, including media when available
func ProcessEvent(event EventStruct, bot *tgbotapi.BotAPI, conf *config.Config) {
	cameraConfig, exists := conf.Cameras[event.Camera]
	if !exists {
		log.Printf("Skipping event from unconfigured camera: %s", event.Camera)
		return
	}

	// Filter by label
	if !contains(cameraConfig.Label, event.Label) {
		log.Printf("Skipping event with unmatched label: %s", event.Label)
		return
	}

	// Filter by score
	if event.Data.TopScore < cameraConfig.Score.MinScore {
		log.Printf("Skipping event with low confidence: %f", event.Data.TopScore)
		return
	}

	// Generate main message using template
	text := generateMessage(event, conf.MessageTemplates)

	// Send main message
	msg := tgbotapi.NewMessage(conf.TelegramChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
	}

	// Send snapshot if available
	if event.HasSnapshot {
		snapshotURL := common.GenerateSnapshotURL(event.ID, conf)  // Moved to common package
		photoMsg := tgbotapi.NewPhotoShare(conf.TelegramChatID, snapshotURL)
		if _, err := bot.Send(photoMsg); err != nil {
			log.Printf("Failed to send snapshot for event %s: %v", event.ID, err)
		}
	}

	// Send video clip if available
	if event.HasClip {
		clipURL := common.GenerateClipURL(event.ID, conf) // Moved to common package
		videoMsg := tgbotapi.NewVideoShare(conf.TelegramChatID, clipURL)
		if _, err := bot.Send(videoMsg); err != nil {
			log.Printf("Failed to send video clip for event %s: %v", event.ID, err)
		}
	}
}

// Helper function to check if a value exists in a slice
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// A function to generate the message based on the template, using the event data
func generateMessage(event EventStruct, templates config.MessageTemplates) string {
	// Create a map of event data to populate the template
	data := map[string]interface{}{
		"camera":   event.Camera,
		"zone":     event.Label, // Assuming label is used as zone here, change if needed
		"label":    event.Label,
		"sublabel": event.SubLabel,
		"score":    event.Data.TopScore,
		"length":   event.EndTime - event.StartTime,
		"severity": event.Data.Type,
		"start_time": time.Unix(int64(event.StartTime), 0).Format(conf.TimeFormat),
		"end_time":   time.Unix(int64(event.EndTime), 0).Format(conf.TimeFormat),
	}

	// Parse the title template
	titleTemplate, err := template.New("title").Parse(templates.TitleTemplate)
	if err != nil {
		log.Printf("Failed to parse title template: %v", err)
		return ""
	}

	// Parse the body template
	bodyTemplate, err := template.New("body").Parse(templates.BodyTemplate)
	if err != nil {
		log.Printf("Failed to parse body template: %v", err)
		return ""
	}

	// Generate the title and body
	var titleBuffer, bodyBuffer bytes.Buffer
	if err := titleTemplate.Execute(&titleBuffer, data); err != nil {
		log.Printf("Failed to execute title template: %v", err)
	}
	if err := bodyTemplate.Execute(&bodyBuffer, data); err != nil {
		log.Printf("Failed to execute body template: %v", err)
	}

	// Return the full message body
	return titleBuffer.String() + "\n" + bodyBuffer.String()
}
