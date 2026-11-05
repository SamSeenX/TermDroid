package bridge

import (
	"strings"
	"testing"
)

func TestBuildScrcpyArgs(t *testing.T) {
	exec := &Executor{
		DeviceID:   "device-1234",
		ScrcpyPath: "/opt/homebrew/bin/scrcpy",
	}

	opts := MirrorOptions{
		TurnScreenOff: true,
		StayAwake:     true,
		Borderless:    true,
		AlwaysOnTop:   true,
		MaxSize:       1024,
		BitRate:       "4M",
		MaxFps:        60,
		VideoCodec:    "h264",
		WindowTitle:   "Custom Window",
	}

	args := exec.BuildScrcpyArgs(opts)
	joined := strings.Join(args, " ")

	expectedTokens := []string{
		"-s device-1234",
		"--turn-screen-off",
		"--stay-awake",
		"--window-borderless",
		"--always-on-top",
		"--max-size=1024",
		"--video-bit-rate=4M",
		"--max-fps=60",
		"--video-codec=h264",
		"--window-title=Custom Window",
	}

	for _, token := range expectedTokens {
		if !strings.Contains(joined, token) {
			t.Errorf("BuildScrcpyArgs output missing expected token: %s. Got: %s", token, joined)
		}
	}
}
