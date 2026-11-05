package ui

import (
	"fmt"
	"strings"
	"termdroid/internal/bridge"

	"github.com/charmbracelet/lipgloss"
)

var colHeaders = []string{
	"Navigation & OS Panels",
	"Volume, Calls & Media",
	"D-Pad Remote & Text",
}

func (m Model) renderKeypadScreen() string {
	var s strings.Builder
	s.WriteString(StyleHeaderMuted.Render("◈ Hardware Keypad & Button Control Deck") + "\n\n")

	cols := bridge.GetKeypadColumns()
	var renderedCols []string

	for colIdx, col := range cols {
		var colBuf strings.Builder
		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1).
			Render(colHeaders[colIdx])
		colBuf.WriteString(header + "\n")

		for rowIdx, act := range col {
			cursor := "  "
			style := StyleNormalRow
			if colIdx == m.KeypadCol && rowIdx == m.KeypadRow {
				cursor = "▶ "
				style = StyleSelectedRow
			}

			tag := act.Shortcut
			if tag == "" {
				tag = act.Code
			}

			line := fmt.Sprintf("%s%-18s [%-5s]", cursor, act.Name, tag)
			colBuf.WriteString(style.Render(line) + "\n")
		}

		renderedCols = append(renderedCols, colBuf.String())
	}

	grid := lipgloss.JoinHorizontal(lipgloss.Top, renderedCols[0], "  ", renderedCols[1], "  ", renderedCols[2])
	s.WriteString(grid)
	s.WriteString("\n\n" + StyleHelp.Render("Arrows: Nav | Enter: Fire | [l] 1-Tap Unlock | [k] Swipe Up | [u] Mute | [h] Home [b] Back [e] Apps [Space] Play"))

	return s.String()
}
