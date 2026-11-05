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

	case "w":
		m.Status = "Negotiating wireless TCP/IP bridge..."
		return m, SwitchWifiCmd(m.Executor)

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

func (m Model) handleUpKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabStorage && m.StorageCursor > 0 {
		m.StorageCursor--
	} else if m.ActiveTab == TabProcesses && m.ProcessCursor > 0 {
		m.ProcessCursor--
	} else if m.ActiveTab == TabActions && m.KeypadRow > 0 {
		m.KeypadRow--
	}
	return m, nil
}

func (m Model) handleDownKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabStorage && m.Storage != nil && m.StorageCursor < len(m.Storage.Items)-1 {
		m.StorageCursor++
	} else if m.ActiveTab == TabProcesses && m.ProcessCursor < len(m.Processes)-1 {
		m.ProcessCursor++
	} else if m.ActiveTab == TabActions {
		cols := bridge.GetKeypadColumns()
		if m.KeypadRow < len(cols[m.KeypadCol])-1 {
			m.KeypadRow++
		}
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
