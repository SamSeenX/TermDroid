package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Devices) == 0 {
		return m.handleWizardKeys(msg)
	}

	if m.PinModalOpen {
		return m.handlePinModalKeys(msg)
	}

	keyStr := msg.String()

	switch keyStr {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.ActiveTab = (m.ActiveTab + 1) % 5
		return m, m.onTabChanged()

	case "shift+tab":
		m.ActiveTab = (m.ActiveTab + 4) % 5
		return m, m.onTabChanged()

	case "1", "2", "3", "4", "5":
		idx := int(msg.Runes[0] - '1')
		m.ActiveTab = Tab(idx)
		return m, m.onTabChanged()

	case "d":
		return m.handleDeviceCycle()

	case "x":
		return m.handleDeviceDisconnect()

	case "a":
		m.Devices = nil
		m.Wizard = InitialWizardState()
		m.Status = "Opening Device Setup Wizard..."
		return m, TickCmd()

	case "r":
		m.IsLoading = true
		return m, m.onTabChanged()

	case "m":
		m.Status = "Launching Scrcpy live mirror..."
		return m, LaunchMirrorCmd(m.Executor, false)

	case "p":
		m.Status = "Launching Scrcpy in Presentation Mode..."
		return m, LaunchMirrorCmd(m.Executor, true)

	case "v":
		if m.ActiveTab == TabMirror {
			m.Status = "Initiating live screen recording..."
			return m, StartRecordingCmd(m.Executor)
		}

	case "c":
		if m.ActiveTab == TabMirror {
			m.Status = "Pushing Mac clipboard to phone..."
			return m, PushClipboardCmd(m.Executor)
		} else if m.ActiveTab == TabStorage {
			m.Status = "Purging application cache directories..."
			return m, PurgeCacheCmd(m.Executor)
		}

	case "y":
		if m.ActiveTab == TabMirror {
			m.Status = "Pulling phone clipboard to Mac..."
			return m, PullClipboardCmd(m.Executor)
		}

	case "s":
		if m.ActiveTab == TabMirror {
			m.Status = "Capturing screen frame..."
			return m, TakeScreenshotCmd(m.Executor)
		}

	case "w":
		m.Status = "Negotiating wireless TCP/IP bridge..."
		return m, SwitchWifiCmd(m.Executor)

	case "l":
		if m.SavedPIN != "" {
			m.Status = "Waking and unlocking device with saved PIN..."
			return m, UnlockDeviceCmd(m.Executor, m.SavedPIN)
		}
		m.PinModalOpen = true
		m.PinInput = ""
		m.Status = "No PIN saved. Enter PIN to enable 1-tap unlock..."
		return m, nil

	case "L":
		m.PinModalOpen = true
		m.PinInput = ""
		m.Status = "Updating saved unlock PIN..."
		return m, nil

	case "k":
		if m.ActiveTab == TabMirror {
			m.Status = "Injecting wake and swipe unlock gesture..."
			return m, UnlockDeviceCmd(m.Executor, "")
		}

	case "left":
		if m.ActiveTab == TabActions && m.KeypadCol > 0 {
			m.KeypadCol--
		}
		return m, nil

	case "right":
		if m.ActiveTab == TabActions && m.KeypadCol < 2 {
			m.KeypadCol++
		}
		return m, nil

	case "up":
		return m.handleUpKey()

	case "down":
		return m.handleDownKey()

	case "enter":
		if m.ActiveTab == TabActions {
			return m.handleKeypadEnter()
		}
		return m, nil

	default:
		if m.ActiveTab == TabActions {
			if act, ok := findActionByShortcut(keyStr); ok {
				m.Status = fmt.Sprintf("Shortcut fired: %s...", act.Name)
				return m, ExecuteKeypadCmd(m.Executor, act)
			}
		}
	}

	return m, nil
}
