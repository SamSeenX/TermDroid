package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderNotificationsTab() string {
	var s strings.Builder
	s.WriteString(StyleHeaderMuted.Render("◈ Live Notifications Stream Deck") + "\n\n")

	if len(m.Notifications) == 0 {
		s.WriteString(StyleHelp.Render("  No active notifications detected. Press [r] to refresh stream.") + "\n")
		return s.String()
	}

	for i, notif := range m.Notifications {
		cursor := "  "
		style := StyleNormalRow
		if i == m.NotifCursor {
			cursor = "▶ "
			style = StyleSelectedRow
		}

		timeBadge := lipgloss.NewStyle().Foreground(ColorMuted).Render(notif.Time)
		pkgBadge := lipgloss.NewStyle().Foreground(ColorSecondary).Render(fmt.Sprintf("[%s]", notif.Package))

		title := notif.Title
		if title == "" {
			title = "(No Title)"
		}
		titleText := lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(title)

		textSnippet := notif.Text
		if len(textSnippet) > 60 {
			textSnippet = textSnippet[:57] + "..."
		}
		if textSnippet == "" {
			textSnippet = "Silent notification"
		}

		line := fmt.Sprintf("%s%s %s %s: %s", cursor, timeBadge, pkgBadge, titleText, textSnippet)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + StyleHelp.Render("▲/▼: Navigate | [r] Refresh Stream | [j] Open Notification Shade on Phone"))
	return s.String()
}
