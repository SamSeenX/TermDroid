package ui

import (
	"fmt"
	"termdroid/internal/bridge"
	"termdroid/internal/diagnostics"

	tea "github.com/charmbracelet/bubbletea"
)

func findActionByShortcut(key string) (bridge.KeyActionItem, bool) {
	for _, col := range bridge.GetKeypadColumns() {
		for _, act := range col {
			if act.Shortcut != "" && (act.Shortcut == key || (act.Shortcut == "Space" && key == " ")) {
				return act, true
			}
		}
	}
	return bridge.KeyActionItem{}, false
}

func (m Model) handleDeviceCycle() (tea.Model, tea.Cmd) {
	if len(m.Devices) > 1 {
		m.ActiveDevice = (m.ActiveDevice + 1) % len(m.Devices)
		m.Executor.DeviceID = m.Devices[m.ActiveDevice].ID
		m.Status = fmt.Sprintf("Switched focus to %s (%s)", m.Devices[m.ActiveDevice].Model, m.Devices[m.ActiveDevice].ID)
		return m, LoadVitalsCmd(m.Executor)
	}
	m.Status = fmt.Sprintf("Current device: %s (No other devices detected)", m.Devices[0].Model)
	return m, nil
}

func (m Model) handleDeviceDisconnect() (tea.Model, tea.Cmd) {
	if len(m.Devices) > 0 {
		dev := m.Devices[m.ActiveDevice]
		if dev.IsWireless {
			m.Status = fmt.Sprintf("Disconnecting %s...", dev.ID)
			return m, DisconnectDeviceCmd(m.Executor, dev.ID)
		}
		m.Status = "USB devices cannot be disconnected via software. Unplug cable."
	}
	return m, nil
}

func (m Model) handleKeypadEnter() (tea.Model, tea.Cmd) {
	cols := bridge.GetKeypadColumns()
	if m.KeypadCol < len(cols) && m.KeypadRow < len(cols[m.KeypadCol]) {
		act := cols[m.KeypadCol][m.KeypadRow]
		m.Status = fmt.Sprintf("Triggering %s...", act.Name)
		return m, ExecuteKeypadCmd(m.Executor, act)
	}
	return m, nil
}

func UnlockDeviceCmd(exec *bridge.Executor, pin string) tea.Cmd {
	return func() tea.Msg {
		err := exec.UnlockDevice(pin)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Unlock failed: %v", err))
		}
		if pin == "" {
			return StatusMsg("Dispatched wake and swipe unlock gesture")
		}
		return StatusMsg("Dispatched wake, swipe, and PIN unlock sequence")
	}
}

func LaunchAudioStreamCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		cmd, err := exec.StartAudioStreamProcess()
		if err != nil {
			return StatusMsg(fmt.Sprintf("Audio forwarder failed: %v", err))
		}
		go cmd.Wait()
		return StatusMsg("Forwarding phone audio to Mac via Opus stream...")
	}
}

func LoadNotificationsCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		items, err := diagnostics.FetchNotifications(exec)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Notification error: %v", err))
		}
		return NotificationsLoadedMsg(items)
	}
}
