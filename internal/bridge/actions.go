package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type KeyAction string

const (
	KeyPower      KeyAction = "26"
	KeyHome       KeyAction = "3"
	KeyBack       KeyAction = "4"
	KeyVolumeUp   KeyAction = "24"
	KeyVolumeDown KeyAction = "25"
	KeyMenu       KeyAction = "82"
)

func (e *Executor) SendKeyEvent(key KeyAction) error {
	_, err := e.RunADBWithTimeout(3*time.Second, "shell", "input", "keyevent", string(key))
	return err
}

func (e *Executor) TakeScreenshot(destPath string) (string, error) {
	if destPath == "" {
		home, _ := os.UserHomeDir()
		timestamp := time.Now().Format("20060102-150405")
		destPath = filepath.Join(home, "Desktop", fmt.Sprintf("termdroid-%s.png", timestamp))
	}

	cmdArgs := []string{}
	if e.DeviceID != "" {
		cmdArgs = append(cmdArgs, "-s", e.DeviceID)
	}
	cmdArgs = append(cmdArgs, "exec-out", "screencap", "-p")

	outFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	cmd := exec.Command(e.ADBPath, cmdArgs...)
	cmd.Stdout = outFile
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func (e *Executor) Tap(x, y int) error {
	args := []string{"shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y)}
	_, err := e.RunADBWithTimeout(3*time.Second, args...)
	return err
}

func (e *Executor) Swipe(x1, y1, x2, y2, durationMs int) error {
	args := []string{
		"shell", "input", "swipe",
		strconv.Itoa(x1), strconv.Itoa(y1),
		strconv.Itoa(x2), strconv.Itoa(y2),
	}
	if durationMs > 0 {
		args = append(args, strconv.Itoa(durationMs))
	}
	_, err := e.RunADBWithTimeout(4*time.Second, args...)
	return err
}

func (e *Executor) OpenURL(url string) error {
	_, err := e.RunADBWithTimeout(4*time.Second, "shell", "am", "start", "-a", "android.intent.action.VIEW", "-d", url)
	return err
}

func (e *Executor) InstallAPK(apkPath string) error {
	_, err := e.RunADBWithTimeout(30*time.Second, "install", "-r", apkPath)
	return err
}

func (e *Executor) ReversePort(port string) error {
	target := fmt.Sprintf("tcp:%s", port)
	_, err := e.RunADBWithTimeout(4*time.Second, "reverse", target, target)
	return err
}

func (e *Executor) RemoveReversePort(port string) error {
	target := fmt.Sprintf("tcp:%s", port)
	_, err := e.RunADBWithTimeout(4*time.Second, "reverse", "--remove", target)
	return err
}

func (e *Executor) PushText(text string) error {
	escaped := fmt.Sprintf("\"%s\"", text)
	_, err := e.RunADBWithTimeout(3*time.Second, "shell", "input", "text", escaped)
	return err
}
