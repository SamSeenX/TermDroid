package bridge

import (
	"fmt"
	"os/exec"
)

type MirrorOptions struct {
	TurnScreenOff    bool
	StayAwake        bool
	Borderless       bool
	RecordPath       string
	MaxSize          int
	BitRate          string
	AlwaysOnTop      bool
	WindowTitle      string
}

func (e *Executor) BuildScrcpyArgs(opts MirrorOptions) []string {
	var args []string
	if e.DeviceID != "" {
		args = append(args, "-s", e.DeviceID)
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
	if opts.MaxSize > 0 {
		args = append(args, fmt.Sprintf("--max-size=%d", opts.MaxSize))
	}
	if opts.BitRate != "" {
		args = append(args, fmt.Sprintf("--video-bit-rate=%s", opts.BitRate))
	}
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
