package bridge

import (
	"fmt"
	"time"
)

type KeyActionItem struct {
	Name     string
	Code     string
	Shortcut string
	IsShell  bool
	CmdArgs  []string
}

func GetKeypadColumns() [][]KeyActionItem {
	col0 := []KeyActionItem{
		{Name: "Home", Code: "3", Shortcut: "h"},
		{Name: "Back", Code: "4", Shortcut: "b"},
		{Name: "Recents / Apps", Code: "187", Shortcut: "e"},
		{Name: "Power Button", Code: "26", Shortcut: "o"},
		{Name: "Screen Sleep", Code: "223"},
		{Name: "Wakeup Device", Code: "224"},
		{Name: "Unlock / Swipe Up", IsShell: true, CmdArgs: []string{"input", "swipe", "500", "1500", "500", "500", "150"}, Shortcut: "k"},
		{Name: "Voice Assistant", Code: "231"},
		{Name: "Menu Key", Code: "82"},
		{Name: "Notifications", IsShell: true, CmdArgs: []string{"cmd", "statusbar", "expand-notifications"}, Shortcut: "j"},
		{Name: "Quick Settings", IsShell: true, CmdArgs: []string{"cmd", "statusbar", "expand-settings"}, Shortcut: "g"},
		{Name: "Screenshot", Code: "120", Shortcut: "s"},
	}

	col1 := []KeyActionItem{
		{Name: "Media Vol (+)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "3", "--adj", "raise", "--show"}, Shortcut: "+"},
		{Name: "Media Vol (-)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "3", "--adj", "lower", "--show"}, Shortcut: "-"},
		{Name: "Call Vol (+)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "0", "--adj", "raise", "--show"}, Shortcut: "c"},
		{Name: "Call Vol (-)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "0", "--adj", "lower", "--show"}, Shortcut: "v"},
		{Name: "Mute Volume", Code: "164", Shortcut: "u"},
		{Name: "Ring Vol (+)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "2", "--adj", "raise", "--show"}},
		{Name: "Ring Vol (-)", IsShell: true, CmdArgs: []string{"cmd", "media_session", "volume", "--stream", "2", "--adj", "lower", "--show"}},
		{Name: "Play / Pause", Code: "85", Shortcut: "Space"},
		{Name: "Next Track", Code: "87", Shortcut: "n"},
		{Name: "Previous Track", Code: "88", Shortcut: "p"},
		{Name: "Stop Playback", Code: "86"},
	}

	col2 := []KeyActionItem{
		{Name: "D-Pad Up", Code: "19"},
		{Name: "D-Pad Down", Code: "20"},
		{Name: "D-Pad Left", Code: "21"},
		{Name: "D-Pad Right", Code: "22"},
		{Name: "D-Pad Center (OK)", Code: "23"},
		{Name: "Enter / Submit", Code: "66"},
		{Name: "Escape", Code: "111"},
		{Name: "Tab Key", Code: "61"},
		{Name: "Backspace", Code: "67"},
		{Name: "Select All", IsShell: true, CmdArgs: []string{"input", "keyevent", "--meta", "113", "29"}},
		{Name: "Copy", IsShell: true, CmdArgs: []string{"input", "keyevent", "--meta", "113", "31"}},
	}

	return [][]KeyActionItem{col0, col1, col2}
}

func (e *Executor) ExecuteKeyAction(action KeyActionItem) (string, error) {
	if action.IsShell {
		args := append([]string{"shell"}, action.CmdArgs...)
		out, err := e.RunADBWithTimeout(3*time.Second, args...)
		if err != nil {
			return "", fmt.Errorf("action %s failed: %w", action.Name, err)
		}
		if out != "" {
			return fmt.Sprintf("Triggered %s (%s)", action.Name, out), nil
		}
		return fmt.Sprintf("Triggered %s", action.Name), nil
	}

	_, err := e.RunADBWithTimeout(3*time.Second, "shell", "input", "keyevent", action.Code)
	if err != nil {
		return "", fmt.Errorf("keyevent %s failed: %w", action.Name, err)
	}

	return fmt.Sprintf("Injected %s (code %s)", action.Name, action.Code), nil
}
