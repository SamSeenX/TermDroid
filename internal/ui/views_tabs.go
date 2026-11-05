package ui

import (
	"fmt"
	"strings"
)

func (m Model) renderActiveTab() string {
	switch m.ActiveTab {
	case TabMirror:
		return m.renderMirrorTab()
	case TabStorage:
		return m.renderStorageTab()
	case TabProcesses:
		return m.renderProcessesTab()
	case TabLogcat:
		return m.renderLogcatTab()
	case TabActions:
		return m.renderActionsTab()
	}
	return ""
}

func (m Model) renderMirrorTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Live Mirror & Desktop Control"))
	s.WriteString("\n\n")
	s.WriteString(StyleNormalRow.Render("  [m] Launch Interactive Scrcpy Stream (Low Latency)\n"))
	s.WriteString(StyleNormalRow.Render("  [p] Launch Presentation Mode (Device Screen Off)\n"))
	s.WriteString(StyleNormalRow.Render("  [s] Take Instant Screenshot to Desktop\n"))
	s.WriteString(StyleNormalRow.Render("  [w] Switch from USB to Wireless TCP/IP (Port 5555)\n\n"))

	if len(m.Devices) > 0 {
		dev := m.Devices[m.ActiveDevice]
		mode := "USB Debugging"
		if dev.IsWireless {
			mode = "Wireless TCP/IP"
		}
		s.WriteString(StyleHelp.Render(fmt.Sprintf("  Target Device: %s (%s) via %s\n", dev.Model, dev.ID, mode)))
	} else {
		s.WriteString(StyleErrorBadge.Render("  No Android device detected. Plug in USB cable.\n"))
	}
	return s.String()
}

func (m Model) renderStorageTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("DiskMan Mobile Storage Breakdown (/sdcard)"))
	s.WriteString("\n\n")

	if m.Storage == nil || len(m.Storage.Items) == 0 {
		s.WriteString(StyleHelp.Render("  No storage data. Press [r] to trigger deep scan.\n"))
		return s.String()
	}

	header := fmt.Sprintf("  Used: %s | Free: %s | Total: %s (%s utilized)\n\n",
		m.Storage.UsedSpace, m.Storage.FreeSpace, m.Storage.TotalSpace, m.Storage.Percent)
	s.WriteString(StyleHelp.Render(header))

	limit := 10
	if len(m.Storage.Items) < limit {
		limit = len(m.Storage.Items)
	}

	for i := 0; i < limit; i++ {
		item := m.Storage.Items[i]
		cursor := "  "
		style := StyleNormalRow
		if i == m.StorageCursor {
			cursor = "> "
			style = StyleSelectedRow
		}
		line := fmt.Sprintf("%s%-10s  %-20s  %s", cursor, item.SizeHuman, item.Category, item.Path)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + StyleHelp.Render("  Use j/k to navigate items. Press [r] to re-scan."))
	return s.String()
}

func (m Model) renderProcessesTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Active Android Processes (ps -A)"))
	s.WriteString("\n\n")

	if len(m.Processes) == 0 {
		s.WriteString(StyleHelp.Render("  No processes found. Press [r] to refresh.\n"))
		return s.String()
	}

	limit := 12
	if len(m.Processes) < limit {
		limit = len(m.Processes)
	}

	s.WriteString(StyleHelp.Render("  PID     USER         NAME\n"))
	for i := 0; i < limit; i++ {
		proc := m.Processes[i]
		cursor := " "
		style := StyleNormalRow
		if i == m.ProcessCursor {
			cursor = ">"
			style = StyleSelectedRow
		}
		line := fmt.Sprintf("%s %-7s %-12s %s", cursor, proc.PID, proc.User, proc.Name)
		s.WriteString(style.Render(line) + "\n")
	}
	return s.String()
}

func (m Model) renderLogcatTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Real-Time System Logcat Viewport"))
	s.WriteString("\n\n")

	if len(m.LogLines) == 0 {
		s.WriteString(StyleHelp.Render("  No log records captured. Press [r] to pull.\n"))
		return s.String()
	}

	start := 0
	if len(m.LogLines) > 12 {
		start = len(m.LogLines) - 12
	}

	for _, line := range m.LogLines[start:] {
		if len(line) > 90 {
			line = line[:87] + "..."
		}
		s.WriteString("  " + StyleNormalRow.Render(line) + "\n")
	}
	return s.String()
}

func (m Model) renderActionsTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Hardware & System Simulation Shortcuts"))
	s.WriteString("\n\n")
	s.WriteString(StyleNormalRow.Render("  [h] Send Home Keyevent (3)\n"))
	s.WriteString(StyleNormalRow.Render("  [b] Send Back Keyevent (4)\n"))
	s.WriteString(StyleNormalRow.Render("  [o] Toggle Power Button (26)\n"))
	s.WriteString(StyleNormalRow.Render("  [s] Capture Instant Screen Frame\n"))
	s.WriteString(StyleNormalRow.Render("  [w] Switch to TCP/IP Wireless Mode\n"))
	return s.String()
}
