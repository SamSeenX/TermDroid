package bridge

import (
	"fmt"
	"time"
)

type KeyActionItem struct {
	Category string
	Name     string
	Code     string
	IsShell  bool
	CmdArgs  []string
}

func GetAllKeyActions() []KeyActionItem {
	return []KeyActionItem{
		// Navigation & System
		{Category: "Navigation & System", Name: "Home", Code: "3"},
		{Category: "Navigation & System", Name: "Back", Code: "4"},
		{Category: "Navigation & System", Name: "Recents / App Switcher", Code: "187"},
		{Category: "Navigation & System", Name: "Menu", Code: "82"},
		{Category: "Navigation & System", Name: "Power Button Toggle", Code: "26"},
		{Category: "Navigation & System", Name: "Screen Sleep / Lock", Code: "223"},
		{Category: "Navigation & System", Name: "Wake Up Device", Code: "224"},
		{Category: "Navigation & System", Name: "Voice Assistant", Code: "231"},

		// Audio & Volume
		{Category: "Volume & Audio", Name: "Volume Up (+)", Code: "24"},
		{Category: "Volume & Audio", Name: "Volume Down (-)", Code: "25"},
		{Category: "Volume & Audio", Name: "Volume Mute Toggle", Code: "164"},
		{Category: "Volume & Audio", Name: "Microphone Mute", Code: "91"},

		// Media Playback
		{Category: "Media Playback", Name: "Play / Pause Toggle", Code: "85"},
		{Category: "Media Playback", Name: "Next Track", Code: "87"},
		{Category: "Media Playback", Name: "Previous Track", Code: "88"},
		{Category: "Media Playback", Name: "Stop Playback", Code: "86"},
		{Category: "Media Playback", Name: "Fast Forward", Code: "90"},
		{Category: "Media Playback", Name: "Rewind", Code: "89"},

		// Status Panels & OS Panels
		{Category: "Status Panels", Name: "Expand Notifications", IsShell: true, CmdArgs: []string{"cmd", "statusbar", "expand-notifications"}},
		{Category: "Status Panels", Name: "Expand Quick Settings", IsShell: true, CmdArgs: []string{"cmd", "statusbar", "expand-settings"}},
		{Category: "Status Panels", Name: "Collapse Status Panels", IsShell: true, CmdArgs: []string{"cmd", "statusbar", "collapse"}},
		{Category: "Status Panels", Name: "Hardware Screenshot", Code: "120"},

		// D-Pad & TV Navigation
		{Category: "D-Pad & Remote", Name: "D-Pad Up", Code: "19"},
		{Category: "D-Pad & Remote", Name: "D-Pad Down", Code: "20"},
		{Category: "D-Pad & Remote", Name: "D-Pad Left", Code: "21"},
		{Category: "D-Pad & Remote", Name: "D-Pad Right", Code: "22"},
		{Category: "D-Pad & Remote", Name: "D-Pad Center (OK)", Code: "23"},

		// Text & Inputs
		{Category: "Text & Inputs", Name: "Enter / Submit", Code: "66"},
		{Category: "Text & Inputs", Name: "Escape", Code: "111"},
		{Category: "Text & Inputs", Name: "Tab", Code: "61"},
		{Category: "Text & Inputs", Name: "Backspace / Delete", Code: "67"},
		{Category: "Text & Inputs", Name: "Select All (Ctrl+A)", IsShell: true, CmdArgs: []string{"input", "keyevent", "--meta", "113", "29"}},
		{Category: "Text & Inputs", Name: "Copy (Ctrl+C)", IsShell: true, CmdArgs: []string{"input", "keyevent", "--meta", "113", "31"}},
		{Category: "Text & Inputs", Name: "Paste (Ctrl+V)", IsShell: true, CmdArgs: []string{"input", "keyevent", "--meta", "113", "50"}},
	}
}

func (e *Executor) ExecuteKeyAction(action KeyActionItem) (string, error) {
	if action.IsShell {
		args := append([]string{"shell"}, action.CmdArgs...)
		out, err := e.RunADBWithTimeout(3*time.Second, args...)
		if err != nil {
			return "", fmt.Errorf("command %s failed: %w", action.Name, err)
		}
		return fmt.Sprintf("Executed %s %s", action.Name, out), nil
	}

	_, err := e.RunADBWithTimeout(3*time.Second, "shell", "input", "keyevent", action.Code)
	if err != nil {
		return "", fmt.Errorf("keyevent %s (%s) failed: %w", action.Name, action.Code, err)
	}

	return fmt.Sprintf("Injected %s (keycode %s)", action.Name, action.Code), nil
}
