package ui

import (
	"fmt"
	"strings"
	"termdroid/internal/bridge"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handlePinModalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	switch keyStr {
	case "esc":
		m.PinModalOpen = false
		m.PinInput = ""
		m.Status = "PIN prompt cancelled"
		return m, nil

	case "enter":
		pin := strings.TrimSpace(m.PinInput)
		m.PinModalOpen = false
		m.PinInput = ""
		m.SavedPIN = pin
		bridge.SaveConfig(bridge.Config{SavedPIN: pin})

		if pin != "" {
			m.Status = "PIN saved to ~/.config/termdroid/config.json. Injecting unlock..."
		} else {
			m.Status = "PIN cleared. Performing swipe gesture..."
		}
		return m, UnlockDeviceCmd(m.Executor, pin)

	case "backspace":
		if len(m.PinInput) > 0 {
			m.PinInput = m.PinInput[:len(m.PinInput)-1]
		}
		return m, nil

	default:
		if len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				m.PinInput += string(r)
			}
		}
		return m, nil
	}
}

func (m Model) renderPinModal() string {
	var b strings.Builder
	b.WriteString(StyleHeader.Render("Configure 1-Tap Unlock PIN") + "\n\n")
	b.WriteString("SurfaceFlinger blacks out the lockscreen keypad on video streams.\n")
	b.WriteString("Enter your phone PIN below to save it locally for 1-tap unlock.\n")
	b.WriteString("Pressing [l] will wake, swipe up, inject the PIN, and unlock in <200ms.\n\n")

	masked := strings.Repeat("*", len(m.PinInput))
	if masked == "" {
		if m.SavedPIN != "" {
			masked = fmt.Sprintf("(Saved: %s - type new digits to change)", strings.Repeat("*", len(m.SavedPIN)))
		} else {
			masked = "(type digits or leave blank to clear)"
		}
	}

	box := fmt.Sprintf("PIN Code: [ %s ]", masked)
	b.WriteString(StyleCard.Render(box) + "\n\n")
	b.WriteString(StyleHelp.Render("[Enter] Save & Unlock | [Esc] Cancel"))

	return StyleModal.Render(b.String())
}
