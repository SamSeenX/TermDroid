package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m Model) renderWizard() string {
	var b strings.Builder

	header := StyleHeader.Render("TERMDROID  Device Connection Wizard")
	b.WriteString(StyleCard.Render(header) + "\n\n")

	var btnUSB, btnWifi string
	if m.Wizard.Mode == WizardUSB {
		btnUSB = StyleTabActive.Render("[1] Connect via USB (Auto-Detect)")
		btnWifi = StyleTabInactive.Render("[2] Connect via Wireless Wi-Fi")
	} else {
		btnUSB = StyleTabInactive.Render("[1] Connect via USB (Auto-Detect)")
		btnWifi = StyleTabActive.Render("[2] Connect via Wireless Wi-Fi")
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, btnUSB, "  ", btnWifi) + "\n\n")

	if m.Wizard.Mode == WizardUSB {
		b.WriteString(m.renderUSBChecklist())
		b.WriteString("\n" + StyleNormalRow.Render("Status: "+m.Status) + "\n")
		b.WriteString(StyleHelp.Render("Press [2] or [Tab] Wireless Mode | [r] Force Scan | [q] Exit"))
	} else {
		b.WriteString(m.renderWirelessGuide())
		b.WriteString("\n" + StyleNormalRow.Render("Status: "+m.Status) + "\n")
		b.WriteString(StyleHelp.Render("[Tab] Switch Field | [Enter] Pair & Connect | [Esc] Back to USB | [Ctrl+C] Exit"))
	}

	return b.String()
}

func (m Model) renderUSBChecklist() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("USB Connection Checklist") + "\n\n")
	s.WriteString(StyleNormalRow.Render("  [1] Connect phone to your Mac using a USB-C or USB-A cable.") + "\n")
	s.WriteString(StyleNormalRow.Render("  [2] Unlock your phone and open Settings > Developer Options.") + "\n")
	s.WriteString(StyleNormalRow.Render("  [3] Toggle 'USB Debugging' to ON.") + "\n")
	s.WriteString(StyleNormalRow.Render("  [4] Tap 'Allow USB debugging from this computer' on your phone.") + "\n\n")

	frame := spinnerFrames[m.Wizard.TickCount%len(spinnerFrames)]
	status := fmt.Sprintf("  %s Monitoring for USB connection in background... (poll count: %d)", frame, m.Wizard.TickCount)
	s.WriteString(StyleStatusBadge.Render(status) + "\n\n")
	s.WriteString(StyleHelp.Render("  Once authorized, TermDroid will automatically transition to the cockpit.") + "\n")

	return s.String()
}

func (m Model) renderWirelessGuide() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Android 11+ Wireless Debugging Setup") + "\n\n")
	s.WriteString(StyleNormalRow.Render("  [1] Ensure phone and Mac are connected to the SAME Wi-Fi network.") + "\n")
	s.WriteString(StyleNormalRow.Render("  [2] On phone: Settings > Developer Options > Wireless Debugging (Toggle ON).") + "\n")
	s.WriteString(StyleNormalRow.Render("  [3] Tap 'Pair device with pairing code' to view IP:Port and 6-digit Code.") + "\n\n")

	f1Cursor := "  "
	f2Cursor := "  "
	f1Style := StyleNormalRow
	f2Style := StyleNormalRow

	if m.Wizard.ActiveField == 0 {
		f1Cursor = "> "
		f1Style = StyleSelectedRow
	} else {
		f2Cursor = "> "
		f2Style = StyleSelectedRow
	}

	ipVal := m.Wizard.PairIPPort
	if ipVal == "" && m.Wizard.ActiveField != 0 {
		ipVal = "e.g. 192.168.1.50:41235"
	}
	codeVal := m.Wizard.PairCode
	if codeVal == "" && m.Wizard.ActiveField != 1 {
		codeVal = "6-digit code"
	}

	s.WriteString(f1Style.Render(fmt.Sprintf("%sIP & Port:     [%-25s]", f1Cursor, ipVal)) + "\n")
	s.WriteString(f2Style.Render(fmt.Sprintf("%sPairing Code:  [%-25s]", f2Cursor, codeVal)) + "\n\n")

	if m.Wizard.ErrorMsg != "" {
		s.WriteString(StyleErrorBadge.Render("  "+m.Wizard.ErrorMsg) + "\n\n")
	}

	return s.String()
}
