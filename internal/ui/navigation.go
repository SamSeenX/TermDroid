package ui

import (
	"termdroid/internal/bridge"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleUpKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabStorage && m.StorageCursor > 0 {
		m.StorageCursor--
	} else if m.ActiveTab == TabProcesses && m.ProcessCursor > 0 {
		m.ProcessCursor--
	} else if m.ActiveTab == TabActions && m.KeypadRow > 0 {
		m.KeypadRow--
	} else if m.ActiveTab == TabNotifications && m.NotifCursor > 0 {
		m.NotifCursor--
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
	} else if m.ActiveTab == TabNotifications && m.NotifCursor < len(m.Notifications)-1 {
		m.NotifCursor++
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
	case TabNotifications:
		m.IsLoading = true
		return LoadNotificationsCmd(m.Executor)
	}
	return nil
}
