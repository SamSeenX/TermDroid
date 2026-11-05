package diagnostics

import (
	"strings"
	"termdroid/internal/bridge"
	"time"
)

type NotificationItem struct {
	Package string
	Title   string
	Text    string
	Time    string
	Key     string
}

func FetchNotifications(exec *bridge.Executor) ([]NotificationItem, error) {
	out, err := exec.RunADBWithTimeout(4*time.Second, "shell", "dumpsys", "notification", "--noredact")
	if err != nil {
		return nil, err
	}

	var items []NotificationItem
	lines := strings.Split(out, "\n")
	var current *NotificationItem

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "NotificationRecord(") {
			if current != nil && current.Package != "" {
				items = append(items, *current)
				if len(items) >= 20 {
					break
				}
			}

			pkg := ""
			parts := strings.Split(trimmed, "pkg=")
			if len(parts) > 1 {
				pkgParts := strings.Fields(parts[1])
				if len(pkgParts) > 0 {
					pkg = pkgParts[0]
				}
			}

			current = &NotificationItem{
				Package: pkg,
				Title:   "",
				Text:    "",
				Time:    time.Now().Format("15:04"),
			}
		} else if current != nil {
			if strings.HasPrefix(trimmed, "android.title=") {
				current.Title = strings.TrimPrefix(trimmed, "android.title=")
			} else if strings.HasPrefix(trimmed, "android.text=") {
				current.Text = strings.TrimPrefix(trimmed, "android.text=")
			}
		}
	}

	if current != nil && current.Package != "" && len(items) < 20 {
		items = append(items, *current)
	}

	return items, nil
}
