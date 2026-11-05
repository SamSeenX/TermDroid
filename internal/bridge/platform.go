package bridge

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func ReadHostClipboard() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbpaste")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return strings.TrimSpace(out.String()), nil

	case "windows":
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-Clipboard")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "", err
		}
		return strings.TrimSpace(out.String()), nil

	default:
		if path, err := exec.LookPath("wl-paste"); err == nil {
			cmd := exec.Command(path)
			var out bytes.Buffer
			cmd.Stdout = &out
			if cmd.Run() == nil {
				return strings.TrimSpace(out.String()), nil
			}
		}
		if path, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command(path, "-selection", "clipboard", "-o")
			var out bytes.Buffer
			cmd.Stdout = &out
			if cmd.Run() == nil {
				return strings.TrimSpace(out.String()), nil
			}
		}
		return "", errors.New("neither wl-paste nor xclip found on Linux")
	}
}

func WriteHostClipboard(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()

	case "windows":
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Set-Clipboard", "-Value", text)
		return cmd.Run()

	default:
		if path, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command(path)
			cmd.Stdin = strings.NewReader(text)
			if cmd.Run() == nil {
				return nil
			}
		}
		if path, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command(path, "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(text)
			if cmd.Run() == nil {
				return nil
			}
		}
		return errors.New("neither wl-copy nor xclip found on Linux")
	}
}

func OpenFileWithDefaultViewer(filePath string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-a", "Preview", filePath).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", filePath).Start()
	default:
		return exec.Command("xdg-open", filePath).Start()
	}
}
