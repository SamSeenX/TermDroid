package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary   = lipgloss.Color("#00D7AF")
	ColorSecondary = lipgloss.Color("#5F87AF")
	ColorAccent    = lipgloss.Color("#FF5F87")
	ColorMuted     = lipgloss.Color("#626262")
	ColorBackground= lipgloss.Color("#1A1B26")
	ColorText      = lipgloss.Color("#D0D0D0")
	ColorSuccess   = lipgloss.Color("#5FD787")
	ColorWarning   = lipgloss.Color("#FFAF00")

	StyleHeader = lipgloss.NewStyle().
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
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(0, 1)

	StyleSelectedRow = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#262936"))

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
)
