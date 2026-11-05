package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleTextModalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	switch keyStr {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.TextModalOpen = false
		m.TextInput = ""
		m.Status = "Text injection cancelled"
		return m, nil

	case "enter":
		text := m.TextInput
		m.TextModalOpen = false
		m.TextInput = ""
		if text == "" {
			m.Status = "No text provided to inject"
			return m, nil
		}
		m.Status = fmt.Sprintf("Injecting text to device cursor (%d chars)...", len(text))
		return m, PushTextToDeviceCmd(m.Executor, text)

	case "backspace":
		if len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.TextInput += string(msg.Runes)
		}
		return m, nil
	}
}

func (m Model) renderTextModal() string {
	var b strings.Builder
	b.WriteString(StyleHeader.Render("Inject Text to Phone Cursor") + "\n\n")
	b.WriteString("Type or paste any text below. When you hit Enter, TermDroid will inject\n")
	b.WriteString("the full string directly into the currently focused input field on your phone.\n\n")

	display := m.TextInput
	if display == "" {
		display = "(paste or type text here)"
	}

	box := fmt.Sprintf("Text: [ %s ]", display)
	b.WriteString(StyleCard.Render(box) + "\n\n")
	b.WriteString(StyleHelp.Render("[Enter] Inject Text to Device | [Esc] Cancel"))

	return StyleModal.Render(b.String())
}
