package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
		return m.renderKeypadScreen()
	case TabNotifications:
		return m.renderNotificationsTab()
	}
	return ""
}

func (m Model) renderMirrorTab() string {
	var s strings.Builder
	s.WriteString(StyleHeaderMuted.Render("◈ Live Mirror & Desktop Stream Control") + "\n\n")

	hotkey := func(k string, desc string) string {
		tag := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(ColorSecondary).Render(" " + k + " ")
		text := lipgloss.NewStyle().Foreground(ColorText).Render(" " + desc)
		return fmt.Sprintf("  %s%s\n", tag, text)
	}

	s.WriteString(hotkey("m", "Launch Interactive Scrcpy Stream (Low Latency, 1200p, 60fps)"))
	s.WriteString(hotkey("p", "Launch Presentation Mode (Device Display Turned Off)"))
	s.WriteString(hotkey("v", "Start Real-Time Screen Recording (Direct MP4 to Desktop)"))
	s.WriteString(hotkey("o", "Stream Mobile Audio to Mac Speakers (Headless Opus Audio)"))
	s.WriteString(hotkey("s", "Capture High-Res Screenshot to Desktop"))
	s.WriteString(hotkey("l", "1-Tap Instant Lockscreen Unlock (Auto-inject PIN from sleep)"))
	s.WriteString(hotkey("t", "Type & Send Text to Device (Inject text to current cursor)"))
	s.WriteString(hotkey("c", "Push Mac Clipboard ⇋ [y] Pull Phone Clipboard"))
	s.WriteString(hotkey("w", "Switch USB Connection to Wireless TCP/IP (Port 5555)"))
	s.WriteString("\n")
	s.WriteString(StyleHelp.Render("  Window Shortcuts: [Opt+h] Home | [Opt+b] Back | [Opt+s] Apps | [Opt+↑/↓] Volume | [Opt+p] Power") + "\n")
	s.WriteString("\n")

	if len(m.Devices) > 0 {
		dev := m.Devices[m.ActiveDevice]
		mode := "USB Bridge"
		badgeColor := ColorPrimary
		if dev.IsWireless {
			mode = "Wireless Wi-Fi"
			badgeColor = ColorSecondary
		}
		badge := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(badgeColor).Render(" " + mode + " ")
		s.WriteString(fmt.Sprintf("  Target: %s [%s] %s\n", dev.Model, dev.ID, badge))
	} else {
		s.WriteString(StyleErrorBadge.Render("  No Android device detected. Connect via USB or Wi-Fi.") + "\n")
	}
	return s.String()
}

func (m Model) renderStorageTab() string {
	var s strings.Builder
	currentPath := m.StoragePath
	if currentPath == "" {
		currentPath = "/sdcard"
	}

	title := fmt.Sprintf("DiskMan Storage: %s", currentPath)
	s.WriteString(StyleHeaderMuted.Render(title) + "\n\n")

	if m.Storage == nil || len(m.Storage.Items) == 0 {
		s.WriteString(StyleHelp.Render("  Scanning directory breakdown... Press [r] to refresh.") + "\n")
		return s.String()
	}

	header := fmt.Sprintf("  Used: %s | Free: %s | Total: %s (%s utilized)",
		m.Storage.UsedSpace, m.Storage.FreeSpace, m.Storage.TotalSpace, m.Storage.Percent)
	s.WriteString(StyleHelp.Render(header) + "\n\n")

	limit := 12
	if len(m.Storage.Items) < limit {
		limit = len(m.Storage.Items)
	}

	for i := 0; i < limit; i++ {
		item := m.Storage.Items[i]
		cursor := "  "
		style := StyleNormalRow
		if i == m.StorageCursor {
			cursor = "▶ "
			style = StyleSelectedRow
		}
		icon := "📁 "
		if !item.IsDir {
			icon = "📄 "
		}
		line := fmt.Sprintf("%s%-10s  %-16s  %s%s", cursor, item.SizeHuman, item.Category, icon, item.Path)
		s.WriteString(style.Render(line) + "\n")
	}

	navHint := "▲/▼: Select | [Enter]: Open | [Esc/Backspace]: Go Up | [x]: Delete | [c]: Purge Cache"
	if len(m.StorageStack) > 0 {
		navHint = fmt.Sprintf("▲/▼: Select | [Enter]: Open | [Esc/Backspace]: Back | [x]: Delete | [r]: Re-scan")
	}
	if m.ConfirmDelete {
		navHint = fmt.Sprintf("CONFIRM DELETION of %s? Press [y] to delete permanently, [Esc] to cancel", m.DeleteTarget)
		s.WriteString("\n" + StyleErrorBadge.Render("  "+navHint))
	} else {
		s.WriteString("\n" + StyleHelp.Render("  "+navHint))
	}
	return s.String()
}

func (m Model) renderProcessesTab() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Active Android Processes (ps -A)") + "\n\n")

	if len(m.Processes) == 0 {
		s.WriteString(StyleHelp.Render("  No processes found. Press [r] to refresh.") + "\n")
		return s.String()
	}

	limit := 12
	if len(m.Processes) < limit {
		limit = len(m.Processes)
	}

	s.WriteString(StyleHelp.Render("  PID     USER         NAME") + "\n")
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
	s.WriteString(StyleHeader.Render("Real-Time System Logcat Viewport") + "\n\n")

	if len(m.LogLines) == 0 {
		s.WriteString(StyleHelp.Render("  No log records captured. Press [r] to pull.") + "\n")
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
