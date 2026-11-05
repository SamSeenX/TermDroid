package storage

import (
	"fmt"
	"strings"
)

type StorageItem struct {
	Path        string
	SizeHuman   string
	SizeBytes   int64
	IsDir       bool
	Category    string
	PercentUsed float64
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func CategorizePath(path string) string {
	lower := strings.ToLower(path)
	switch {
	case containsSub(lower, "dcim", "camera", "pictures", "screenshots"):
		return "Photos & Media"
	case containsSub(lower, "whatsapp", "telegram", "viber"):
		return "Messaging Media"
	case containsSub(lower, "download", "downloads"):
		return "Downloads"
	case containsSub(lower, "music", "podcasts", "audio"):
		return "Audio & Music"
	case containsSub(lower, "movies", "video"):
		return "Video"
	case containsSub(lower, "cache", "temp"):
		return "App Cache"
	default:
		return "Other Storage"
	}
}

func containsSub(s string, sub ...string) bool {
	for _, item := range sub {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}
