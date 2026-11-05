package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const AsciiLogoRaw = `
 _____ _____ ____  __  __ ____  ____   ___ ___ ____  
|_   _| ____|  _ \|  \/  |  _ \|  _ \ / _ \_ _|  _ \ 
  | | |  _| | |_) | |\/| | | | | |_) | | | | || | | |
  | | | |___|  _ <| |  | | |_| |  _ <| |_| | || |_| |
  |_| |_____|_| \_\_|  |_|____/|_| \_\\___/|___|____/ 
`

func renderAsciiHeader() string {
	trimmed := strings.Trim(AsciiLogoRaw, "\n")
	logoStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary)

	return logoStyle.Render(trimmed)
}
