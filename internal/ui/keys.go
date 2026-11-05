package ui

import (
	"fmt"
	"termdroid/internal/bridge"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Devices) == 0 {
		return m.handleWizardKeys(msg)
	}

	switch msg.String() {
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
		if len(m.Devices) > 1 {
			m.ActiveDevice = (m.ActiveDevice + 1) % len(m.Devices)
			m.Executor.DeviceID = m.Devices[m.ActiveDevice].ID
			m.Status = fmt.Sprintf("Switched focus to %s (%s)", m.Devices[m.ActiveDevice].Model, m.Devices[m.ActiveDevice].ID)
			return m, LoadVitalsCmd(m.Executor)
		}
		m.Status = fmt.Sprintf("Current device: %s (No other devices detected)", m.Devices[0].Model)
		return m, nil

	case "x":
		if len(m.Devices) > 0 {
			dev := m.Devices[m.ActiveDevice]
			if dev.IsWireless {
				m.Status = fmt.Sprintf("Disconnecting %s...", dev.ID)
				return m, DisconnectDeviceCmd(m.Executor, dev.ID)
			}
			m.Status = "USB devices cannot be disconnected via software. Unplug cable."
			return m, nil
		}

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

	case "s":
		m.Status = "Capturing screen frame..."
		return m, TakeScreenshotCmd(m.Executor)

	case "w":
		m.Status = "Negotiating wireless TCP/IP bridge..."
		return m, SwitchWifiCmd(m.Executor)

	case "up", "k":
		return m.handleUpKey()

	case "down", "j":
		return m.handleDownKey()

	case "enter":
		if m.ActiveTab == TabActions {
			actions := bridge.GetAllKeyActions()
			if m.KeypadCursor >= 0 && m.KeypadCursor < len(actions) {
				act := actions[m.KeypadCursor]
				m.Status = fmt.Sprintf("Triggering %s...", act.Name)
				return m, ExecuteKeypadCmd(m.Executor, act)
			}
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleUpKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabStorage && m.StorageCursor > 0 {
		m.StorageCursor--
	} else if m.ActiveTab == TabProcesses && m.ProcessCursor > 0 {
		m.ProcessCursor--
	} else if m.ActiveTab == TabActions {
		actions := bridge.GetAllKeyActions()
		if m.KeypadCursor > 0 {
			m.KeypadCursor--
		} else {
			m.KeypadCursor = len(actions) - 1
		}
	}
	return m, nil
}

func (m Model) handleDownKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabStorage && m.Storage != nil && m.StorageCursor < len(m.Storage.Items)-1 {
		m.StorageCursor++
	} else if m.ActiveTab == TabProcesses && m.ProcessCursor < len(m.Processes)-1 {
		m.ProcessCursor++
	} else if m.ActiveTab == TabActions {
		actions := bridge.GetAllKeyActions()
		m.KeypadCursor = (m.KeypadCursor + 1) % len(actions)
	}
	return m, nil
}

func (m *Model) onTabChanged() tea.Cmd {
	if len(m.Devices) == 0 {
		return LoadDevicesCmd(m.Executor)
	}

	switch m.ActiveTab {
	case TabMirror:
		return LoadVitalsCmd(m.Executor)
	case TabStorage:
		if m.Storage == nil {
			m.IsLoading = true
			return LoadStorageCmd(m.Executor)
		}
	case TabProcesses:
		m.IsLoading = true
		return LoadProcessesCmd(m.Executor)
	case TabLogcat:
		m.IsLoading = true
		return LoadLogcatCmd(m.Executor)
	case TabActions:
		return nil
	}
	return nil
}
