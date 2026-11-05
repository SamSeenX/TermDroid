package bridge

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func (e *Executor) PushClipboardToDevice() (string, error) {
	pbCmd := exec.Command("pbpaste")
	var pbOut bytes.Buffer
	pbCmd.Stdout = &pbOut
	if err := pbCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read mac clipboard: %w", err)
	}

	text := strings.TrimSpace(pbOut.String())
	if text == "" {
		return "", errors.New("mac clipboard is empty")
	}

	_, err := e.RunADBWithTimeout(4*time.Second, "shell", "cmd", "clipboard", "set", "text", text)
	if err != nil {
		escaped := fmt.Sprintf("\"%s\"", text)
		_, err = e.RunADBWithTimeout(4*time.Second, "shell", "input", "text", escaped)
	}

	if err != nil {
		return "", fmt.Errorf("failed to sync clipboard to phone: %w", err)
	}

	preview := text
	if len(preview) > 25 {
		preview = preview[:22] + "..."
	}
	return fmt.Sprintf("Pushed to phone clipboard: '%s'", preview), nil
}

func (e *Executor) PullClipboardFromDevice() (string, error) {
	out, err := e.RunADBWithTimeout(4*time.Second, "shell", "cmd", "clipboard", "get")
	text := strings.TrimSpace(out)

	if err != nil || text == "" {
		dumpOut, dumpErr := e.RunADBWithTimeout(3*time.Second, "shell", "dumpsys", "clipboard")
		if dumpErr == nil {
			for _, line := range strings.Split(dumpOut, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "text=") {
					val := strings.TrimPrefix(line, "text=")
					if val != "" {
						text = val
						break
					}
				}
			}
		}
	}

	if text == "" {
		return "Tip: Android 13+ restricts background clipboard reads. When Scrcpy mirror is open [m], copied text syncs to Mac automatically.", nil
	}

	pbCmd := exec.Command("pbcopy")
	pbCmd.Stdin = strings.NewReader(text)
	if err := pbCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to write to mac clipboard: %w", err)
	}

	preview := text
	if len(preview) > 25 {
		preview = preview[:22] + "..."
	}
	return fmt.Sprintf("Copied from phone to Mac: '%s'", preview), nil
}
