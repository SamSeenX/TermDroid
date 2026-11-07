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

	if m.TextModalOpen {
		return m.handleTextModalKeys(msg)
	}

	if m.StorageModalOpen {
		return handleStorageModalKey(m, msg)
	}

	keyStr := msg.String()

	if newM, cmd, handled := handleStorageKeys(m, keyStr); handled {
		return newM, cmd
	}

	if newM, cmd, handled := handleMirrorKeys(m, keyStr); handled {
		return newM, cmd
	}

	switch keyStr {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "tab":
		m.ActiveTab = (m.ActiveTab + 1) % 6
		return m, m.onTabChanged()

	case "shift+tab":
		m.ActiveTab = (m.ActiveTab + 5) % 6
		return m, m.onTabChanged()

	case "1", "2", "3", "4", "5", "6":
		idx := int(msg.Runes[0] - '1')
		m.ActiveTab = Tab(idx)
		return m, m.onTabChanged()

	case "d":
		return m.handleDeviceCycle()

	case "x":
		if m.ActiveTab == TabStorage {
			return m.handleDeleteKey()
		}
		return m.handleDeviceDisconnect()

	case "delete":
		if m.ActiveTab == TabStorage {
			return m.handleDeleteKey()
		}

	case "a":
		m.Devices = nil
		m.Wizard = InitialWizardState()
		m.Status = "Opening Device Setup Wizard..."
		return m, TickCmd()

	case "r":
		if m.ActiveTab == TabStorage {
			m.IsScanningStorage = true
			m.IsLoading = true
			m.Status = fmt.Sprintf("Triggering fresh scan of %s...", m.StoragePath)
			return m, LoadStorageCmd(m.Executor, m.StoragePath)
		}
		m.IsLoading = true
		return m, m.onTabChanged()

	case "w":
		m.Status = "Negotiating wireless TCP/IP bridge..."
		return m, SwitchWifiCmd(m.Executor)

	case "t":
		m.TextModalOpen = true
		m.TextInput = ""
		m.Status = "Opening Text Injection prompt..."
		return m, nil

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

	case "y", "Y":
		if m.ConfirmDelete {
			return m.handleConfirmDelete()
		}
		if m.ActiveTab == TabMirror {
			m.Status = "Pulling phone clipboard to Mac..."
			return m, PullClipboardCmd(m.Executor)
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
		return m.handleEnterKey()

	case "backspace", "esc":
		return m.handleBackNavKey()

	default:
		if m.ActiveTab == TabActions {
			return m.handleActionShortcut(keyStr)
		}
	}

	return m, nil
}
