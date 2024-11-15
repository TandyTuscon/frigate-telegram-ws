package frigate

import (
	"log"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/TandyTuscon/frigate-telegram-ws/config"
)

// Generate a formatted message for an event
func generateMessage(event EventStruct, conf *config.Config) string {
	return fmt.Sprintf("Event: %s detected on camera: %s", event.Label, event.Camera)
}

// SaveThumbnail decodes and saves a snapshot image locally
func SaveThumbnail(eventID string, thumbnail string, conf *config.Config) string {
	data, err := base64.StdEncoding.DecodeString(thumbnail)
	if err != nil {
		log.Printf("Failed to decode thumbnail for event %s: %v", eventID, err)
		return ""
	}

	filePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.jpg", eventID))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("Failed to save thumbnail for event %s: %v", eventID, err)
		return ""
	}

	return filePath
}

// SaveClip downloads and saves a video clip locally
func SaveClip(eventID string, conf *config.Config) string {
	url := fmt.Sprintf("%s/api/events/%s/clip.mp4", conf.FrigateURL, eventID)
	filePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.mp4", eventID))

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download clip for event %s: %v", eventID, err)
		return ""
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to save clip for event %s: %v", eventID, err)
		return ""
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		log.Printf("Failed to write clip for event %s: %v", eventID, err)
		return ""
	}

	return filePath
}

// Generate snapshot, clip, and event URLs
func GenerateSnapshotURL(eventID string, conf *config.Config) string {
	return fmt.Sprintf("%s/api/events/%s/snapshot.jpg", conf.FrigateExternalURL, eventID)
}

func GenerateClipURL(eventID string, conf *config.Config) string {
	return fmt.Sprintf("%s/api/events/%s/clip.mp4", conf.FrigateExternalURL, eventID)
}

func GenerateEventURL(camera, label, zone string, conf *config.Config) string {
	return fmt.Sprintf("%s/events?cameras=%s&labels=%s&zones=%s", conf.FrigateExternalURL, camera, label, zone)
}

// ProcessEvent handles an event and sends a message to Telegram
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
	if event.Data["top_score"].(float64) < cameraConfig.Score.MinScore {
		log.Printf("Skipping event with low confidence: %f", event.Data["top_score"])
		return
	}

	// Generate and send the main message
	text := generateMessage(event, conf)
	msg := tgbotapi.NewMessage(conf.TelegramChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
	}

	// Send snapshot if available
	if event.HasSnapshot {
		snapshotPath := GenerateSnapshotURL(event.ID, conf)
		photoMsg := tgbotapi.NewPhoto(conf.TelegramChatID, tgbotapi.FilePath(snapshotPath))
		if _, err := bot.Send(photoMsg); err != nil {
			log.Printf("Failed to send snapshot for event %s: %v", event.ID, err)
		}
	}

	// Send video clip if available
	if event.HasClip {
		clipPath := GenerateClipURL(event.ID, conf)
		videoMsg := tgbotapi.NewVideo(conf.TelegramChatID, tgbotapi.FilePath(clipPath))
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
