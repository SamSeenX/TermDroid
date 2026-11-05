package ui

import (
	"fmt"
	"strings"
	"termdroid/internal/bridge"
)

func (m Model) renderKeypadScreen() string {
	var s strings.Builder
	s.WriteString(StyleHeader.Render("Hardware Keypad & System Simulation Deck") + "\n\n")

	actions := bridge.GetAllKeyActions()
	if len(actions) == 0 {
		return "  No actions defined.\n"
	}

	currentCategory := ""
	for i, act := range actions {
		if act.Category != currentCategory {
			currentCategory = act.Category
			s.WriteString("\n" + StyleHeader.Render("  "+currentCategory) + "\n")
		}

		cursor := "    "
		style := StyleNormalRow
		if i == m.KeypadCursor {
			cursor = "  > "
			style = StyleSelectedRow
		}

		codeBadge := act.Code
		if act.IsShell {
			codeBadge = "shell"
		}

		line := fmt.Sprintf("%s%-30s  [%-5s]", cursor, act.Name, codeBadge)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + StyleHelp.Render("  [j/k or Up/Down] Select Action | [Enter] Fire Button | [d] Cycle Device | [x] Disconnect | [a] Add Device"))
	return s.String()
}
