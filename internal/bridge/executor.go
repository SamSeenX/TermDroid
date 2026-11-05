package bridge

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	ErrDeviceNotFound = errors.New("no android device detected")
	ErrCommandTimeout = errors.New("adb command timed out")
)

type Executor struct {
	ADBPath    string
	ScrcpyPath string
	DeviceID   string
}

func NewExecutor(deviceID string) *Executor {
	adb, err := exec.LookPath("adb")
	if err != nil {
		if runtime.GOOS == "darwin" {
			adb = "/opt/homebrew/bin/adb"
		} else if runtime.GOOS == "windows" {
			adb = "adb.exe"
		} else {
			adb = "/usr/bin/adb"
		}
	}
	scrcpy, err := exec.LookPath("scrcpy")
	if err != nil {
		if runtime.GOOS == "darwin" {
			scrcpy = "/opt/homebrew/bin/scrcpy"
		} else if runtime.GOOS == "windows" {
			scrcpy = "scrcpy.exe"
		} else {
			scrcpy = "/usr/bin/scrcpy"
		}
	}
	return &Executor{
		ADBPath:    adb,
		ScrcpyPath: scrcpy,
		DeviceID:   deviceID,
	}
}

func (e *Executor) RunADB(ctx context.Context, args ...string) (string, error) {
	var finalArgs []string
	if e.DeviceID != "" {
		finalArgs = append(finalArgs, "-s", e.DeviceID)
	}
	finalArgs = append(finalArgs, args...)

	cmd := exec.CommandContext(ctx, e.ADBPath, finalArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr != "" {
			return "", errors.New(errStr)
		}
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (e *Executor) RunADBWithTimeout(timeout time.Duration, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return e.RunADB(ctx, args...)
}
