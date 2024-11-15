package frigate

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "text/template"
    "github.com/TandyTuscon/frigate-telegram-ws/config"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// EventStruct defines the event data
type EventStruct struct {
    ID           string  `json:"id"`
    Camera       string  `json:"camera"`
    Label        string  `json:"label"`
    Score        float64 `json:"score"`
    HasClip      bool    `json:"has_clip"`
    HasSnapshot  bool    `json:"has_snapshot"`
    // Add other fields as needed
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

// GenerateSnapshotURL
func GenerateSnapshotURL(eventID string, conf *config.Config) string {
    return fmt.Sprintf("%s/api/events/%s/snapshot.jpg", conf.FrigateExternalURL, eventID)
}

// GenerateClipURL
func GenerateClipURL(eventID string, conf *config.Config) string {
    return fmt.Sprintf("%s/api/events/%s/clip.mp4", conf.FrigateExternalURL, eventID)
}

// GenerateEventURL
func GenerateEventURL(camera, label, zone string, conf *config.Config) string {
    return fmt.Sprintf("%s/events?cameras=%s&labels=%s&zones=%s", conf.FrigateExternalURL, camera, label, zone)
}

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
    if event.Score < cameraConfig.Score.MinScore {
        log.Printf("Skipping event with low confidence: %f", event.Score)
        return
    }

    // Generate main message
    text := generateMessage(event, conf.MessageTemplates)

    // Send main message
    msg := tgbotapi.NewMessage(conf.TelegramChatID, text)
    msg.ParseMode = tgbotapi.ModeMarkdown
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Failed to send Telegram message: %v", err)
    }

    // Send snapshot if available
    if event.HasSnapshot {
        snapshotURL := GenerateSnapshotURL(event.ID, conf)
        photoMsg := tgbotapi.NewPhotoShare(conf.TelegramChatID, snapshotURL)
        if _, err := bot.Send(photoMsg); err != nil {
            log.Printf("Failed to send snapshot for event %s: %v", event.ID, err)
        }
    }

    // Send video clip if available
    if event.HasClip {
        clipURL := GenerateClipURL(event.ID, conf)
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
