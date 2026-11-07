package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if len(m.Devices) == 0 {
		return m.renderWizard()
	}

	if m.PinModalOpen {
		return m.renderPinModal()
	}

	if m.TextModalOpen {
		return m.renderTextModal()
	}

	if m.StorageModalOpen {
		return m.renderStorageModal()
	}

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
	ascii := renderAsciiHeader()
	subtitle := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render("macOS ⇋ Android Control Cockpit & Diagnostics")
	titleLine := fmt.Sprintf("%s\n  %s", ascii, subtitle)

	var vitals string
	if m.Vitals != nil {
		chargingIcon := "🔋 "
		chargeStatus := "Discharging"
		if m.Vitals.BatteryCharging {
			chargingIcon = "⚡ "
			chargeStatus = fmt.Sprintf("Charging (%s)", m.Vitals.ChargingSource)
		}

		battColor := ColorSuccess
		wifiName := m.Vitals.WifiSSID
		if wifiName == "" {
			wifiName = "Disconnected"
		}

		vitals = fmt.Sprintf("%sBattery: %s [%s]  🌡️ Temp: %s  📶 Wi-Fi: %s",
			chargingIcon,
			lipgloss.NewStyle().Bold(true).Foreground(battColor).Render(m.Vitals.BatteryLevel),
			lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render(chargeStatus),
			m.Vitals.BatteryTemperature,
			lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(wifiName),
		)
	} else {
		vitals = "🔋 Battery: Checking status..."
	}

	hud := StyleCard.Render(fmt.Sprintf("%s\n%s", titleLine, StyleHelp.Render(vitals)))
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
	shortcuts := StyleHelp.Render("Tab: Tabs | l: 1-Tap Unlock | m: Mirror | d: Switch Device | s: Shot | r: Refresh | q: Quit")
	return fmt.Sprintf("\n%s\n%s", status, shortcuts)
}
