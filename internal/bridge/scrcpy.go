package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type MirrorOptions struct {
	TurnScreenOff bool
	StayAwake     bool
	Borderless    bool
	RecordPath    string
	MaxSize       int
	BitRate       string
	MaxFps        int
	VideoCodec    string
	AlwaysOnTop   bool
	WindowTitle   string
	AudioOnly     bool
}

func (e *Executor) BuildScrcpyArgs(opts MirrorOptions) []string {
	var args []string
	if e.DeviceID != "" {
		args = append(args, "-s", e.DeviceID)
	}

	if opts.AudioOnly {
		args = append(args, "--no-video", "--audio-codec=opus", "--audio-buffer=40")
		return args
	}

	if opts.TurnScreenOff {
		args = append(args, "--turn-screen-off")
	}
	if opts.StayAwake {
		args = append(args, "--stay-awake")
	}
	if opts.Borderless {
		args = append(args, "--window-borderless")
	}
	if opts.AlwaysOnTop {
		args = append(args, "--always-on-top")
	}

	maxSize := opts.MaxSize
	if maxSize <= 0 {
		maxSize = 1200
	}
	args = append(args, fmt.Sprintf("--max-size=%d", maxSize))

	bitRate := opts.BitRate
	if bitRate == "" {
		bitRate = "6M"
	}
	args = append(args, fmt.Sprintf("--video-bit-rate=%s", bitRate))

	maxFps := opts.MaxFps
	if maxFps <= 0 {
		maxFps = 60
	}
	args = append(args, fmt.Sprintf("--max-fps=%d", maxFps))

	codec := opts.VideoCodec
	if codec == "" {
		codec = "h264"
	}
	args = append(args, fmt.Sprintf("--video-codec=%s", codec))

	if opts.RecordPath != "" {
		args = append(args, fmt.Sprintf("--record=%s", opts.RecordPath))
	}

	title := opts.WindowTitle
	if title == "" {
		title = "TermDroid Live Stream"
	}
	args = append(args, fmt.Sprintf("--window-title=%s", title))

	return args
}

func (e *Executor) StartMirrorProcess(opts MirrorOptions) (*exec.Cmd, error) {
	args := e.BuildScrcpyArgs(opts)
	cmd := exec.Command(e.ScrcpyPath, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (e *Executor) StartRecordingProcess(destPath string) (*exec.Cmd, string, error) {
	if destPath == "" {
		home, _ := os.UserHomeDir()
		timestamp := time.Now().Format("20060102-150405")
		destPath = filepath.Join(home, "Desktop", fmt.Sprintf("termdroid-recording-%s.mp4", timestamp))
	}

	opts := MirrorOptions{
		StayAwake:   true,
		RecordPath:  destPath,
		MaxSize:     1200,
		BitRate:     "6M",
		MaxFps:      60,
		WindowTitle: "TermDroid Recording Active",
	}

	cmd, err := e.StartMirrorProcess(opts)
	if err != nil {
		return nil, "", err
	}

	return cmd, destPath, nil
}

func (e *Executor) StartAudioStreamProcess() (*exec.Cmd, error) {
	opts := MirrorOptions{
		AudioOnly:   true,
		WindowTitle: "TermDroid Mobile Audio Forwarder",
	}
	return e.StartMirrorProcess(opts)
}
