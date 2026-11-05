package diagnostics

import (
	"fmt"
	"strings"
	"termdroid/internal/bridge"
	"time"
)

type LogcatFilter string

const (
	LogVerbose LogcatFilter = "*:V"
	LogDebug   LogcatFilter = "*:D"
	LogInfo    LogcatFilter = "*:I"
	LogWarning LogcatFilter = "*:W"
	LogError   LogcatFilter = "*:E"
)

func FetchRecentLogs(exec *bridge.Executor, lines int, level LogcatFilter) ([]string, error) {
	if lines <= 0 {
		lines = 60
	}
	if level == "" {
		level = LogInfo
	}

	filterArg := string(level)
	linesArg := fmt.Sprintf("%d", lines)

	out, err := exec.RunADBWithTimeout(4*time.Second, "logcat", "-d", "-v", "time", "-t", linesArg, filterArg)
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(out, "\n")
	var result []string
	for _, l := range rawLines {
		trimmed := strings.TrimRight(l, "\r\n")
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result, nil
}

func ClearLogcat(exec *bridge.Executor) error {
	_, err := exec.RunADBWithTimeout(3*time.Second, "logcat", "-c")
	return err
}
