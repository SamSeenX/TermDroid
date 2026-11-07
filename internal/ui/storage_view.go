package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderStorageTab() string {
	var s strings.Builder
	currentPath := m.StoragePath
	if currentPath == "" {
		currentPath = "/sdcard"
	}

	title := fmt.Sprintf("DiskMan Storage: %s", currentPath)
	s.WriteString(StyleHeaderMuted.Render(title))

	if len(m.SelectedFiles) > 0 {
		selCount := 0
		for _, v := range m.SelectedFiles {
			if v {
				selCount++
			}
		}
		if selCount > 0 {
			selBadge := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(ColorPrimary).
				Render(fmt.Sprintf(" %d MARKED ", selCount))
			s.WriteString("  " + selBadge)
		}
	}

	if len(m.ClipboardPaths) > 0 {
		mode := "COPY"
		bg := ColorSecondary
		if m.ClipboardIsCut {
			mode = "CUT"
			bg = ColorWarning
		}
		clipBadge := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(bg).
			Render(fmt.Sprintf(" %s %d ITEMS ", mode, len(m.ClipboardPaths)))
		s.WriteString("  " + clipBadge)
	}

	s.WriteString("\n\n")

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

		mark := "[ ] "
		if m.SelectedFiles[item.Path] {
			mark = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("[✓] ")
		}

		displayPath := filepath.Base(item.Path)
		if displayPath == "." || displayPath == "/" {
			displayPath = item.Path
		}

		line := fmt.Sprintf("%s%s%-10s  %-15s  %s%s", cursor, mark, item.SizeHuman, item.Category, icon, displayPath)
		s.WriteString(style.Render(line) + "\n")
	}

	if m.ConfirmDelete {
		count := m.getSelectedOrCursorCount()
		prompt := fmt.Sprintf("CONFIRM DELETION of %d item(s)? Press [y] to delete permanently, [Esc] to cancel", count)
		s.WriteString("\n" + StyleErrorBadge.Render("  "+prompt))
		return s.String()
	}

	navHint := "▲/▼: Select | [Space]: Mark | [a]: All | [d]: Pull to PC | [u]: Upload | [c]: Copy | [m]: Cut | [p]: Paste"
	subHint := "[Enter]: Open | [Esc/Bksp]: Back | [n]: Mkdir | [R]: Rename | [x]: Delete | [r]: Rescan"
	s.WriteString("\n" + StyleHelp.Render("  "+navHint) + "\n" + StyleHelp.Render("  "+subHint))
	return s.String()
}

func (m Model) getSelectedOrCursorCount() int {
	selCount := 0
	for _, v := range m.SelectedFiles {
		if v {
			selCount++
		}
	}
	if selCount > 0 {
		return selCount
	}
	return 1
}
