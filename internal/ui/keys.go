package ui

import (
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

	case "r":
		m.IsLoading = true
		return m, m.onTabChanged()

	case "m":
		m.Status = "Launching Scrcpy live mirror..."
		return m, LaunchMirrorCmd(m.Executor, false)

	case "p":
		m.Status = "Launching Scrcpy in Presentation Mode (screen off)..."
		return m, LaunchMirrorCmd(m.Executor, true)

	case "s":
		m.Status = "Capturing screen frame..."
		return m, TakeScreenshotCmd(m.Executor)

	case "w":
		m.Status = "Negotiating wireless TCP/IP bridge..."
		return m, SwitchWifiCmd(m.Executor)

	case "up", "k":
		if m.ActiveTab == TabStorage && m.StorageCursor > 0 {
			m.StorageCursor--
		} else if m.ActiveTab == TabProcesses && m.ProcessCursor > 0 {
			m.ProcessCursor--
		}
		return m, nil

	case "down", "j":
		if m.ActiveTab == TabStorage && m.Storage != nil && m.StorageCursor < len(m.Storage.Items)-1 {
			m.StorageCursor++
		} else if m.ActiveTab == TabProcesses && m.ProcessCursor < len(m.Processes)-1 {
			m.ProcessCursor++
		}
		return m, nil

	case "h":
		return m, SendKeyCmd(m.Executor, bridge.KeyHome)
	case "b":
		return m, SendKeyCmd(m.Executor, bridge.KeyBack)
	case "o":
		return m, SendKeyCmd(m.Executor, bridge.KeyPower)
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
		return LoadVitalsCmd(m.Executor)
	}
	return nil
}
