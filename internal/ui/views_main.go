package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	if m.IsLoading {
		b.WriteString(StyleNormalRow.Render("  [Loading data from device...]"))
		b.WriteString("\n\n")
	} else {
		b.WriteString(m.renderActiveTab())
		b.WriteString("\n")
	}

	b.WriteString(m.renderStatusBar())
	return b.String()
}

func (m Model) renderHeader() string {
	title := StyleHeader.Render("TERMDROID  macOS Android Control Cockpit")
	var vitals string
	if m.Vitals != nil {
		charging := ""
		if m.Vitals.BatteryCharging {
			charging = "⚡ "
		}
		vitals = fmt.Sprintf("Battery: %s%s (%s)  Temp: %s",
			charging,
			m.Vitals.BatteryLevel,
			m.Vitals.BatteryHealth,
			m.Vitals.BatteryTemperature,
		)
	} else {
		vitals = "Battery: Checking..."
	}

	hud := StyleCard.Render(fmt.Sprintf("%s\n%s", title, StyleHelp.Render(vitals)))
	return hud
}

func (m Model) renderTabs() string {
	var renderedTabs []string
	for i, name := range TabNames {
		if Tab(i) == m.ActiveTab {
			renderedTabs = append(renderedTabs, StyleTabActive.Render(name))
		} else {
			renderedTabs = append(renderedTabs, StyleTabInactive.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m Model) renderStatusBar() string {
	status := StyleNormalRow.Render("Status: " + m.Status)
	shortcuts := StyleHelp.Render("Tab: Cycle Tabs | m: Mirror | p: Pres. Mode | s: Shot | w: Wi-Fi | r: Refresh | q: Quit")
	return fmt.Sprintf("\n%s\n%s", status, shortcuts)
}
