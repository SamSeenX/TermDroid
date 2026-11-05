package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary   = lipgloss.Color("#00FFC2")
	ColorSecondary = lipgloss.Color("#7957D5")
	ColorAccent    = lipgloss.Color("#FF3860")
	ColorMuted     = lipgloss.Color("#7A7E9D")
	ColorBackground= lipgloss.Color("#0F111A")
	ColorSurface   = lipgloss.Color("#1A1D2B")
	ColorText      = lipgloss.Color("#ECEFF4")
	ColorSuccess   = lipgloss.Color("#48C774")
	ColorWarning   = lipgloss.Color("#FFDD57")

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorPrimary).
			Padding(0, 2)

	StyleHeaderMuted = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorPrimary).
			Padding(0, 2)

	StyleTabInactive = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 2)

	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1)

	StyleSelectedRow = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#2B2144"))

	StyleNormalRow = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleStatusBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	StyleErrorBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleModal = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)
)
