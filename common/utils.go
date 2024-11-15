
package common

import (
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "github.com/TandyTuscon/frigate-telegram-ws/config"
)

// SaveThumbnail decodes and saves a snapshot image locally
func SaveThumbnail(eventID string, thumbnail string, conf *config.Config) string {
    data, err := base64.StdEncoding.DecodeString(thumbnail)
    if err != nil {
        return ""
    }

    filePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.jpg", eventID))
    if err := os.WriteFile(filePath, data, 0644); err != nil {
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
        return ""
    }
    defer resp.Body.Close()

    file, err := os.Create(filePath)
    if err != nil {
        return ""
    }
    defer file.Close()

    if _, err := io.Copy(file, resp.Body); err != nil {
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
