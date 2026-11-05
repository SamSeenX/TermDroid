package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type WizardMode int

const (
	WizardUSB WizardMode = iota
	WizardWireless
)

type TickMsg time.Time

func TickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type WizardState struct {
	Mode        WizardMode
	TickCount   int
	PairIPPort  string
	PairCode    string
	ActiveField int
	ErrorMsg    string
	SuccessMsg  string
	IsPairing   bool
}

func InitialWizardState() WizardState {
	return WizardState{
		Mode:        WizardUSB,
		TickCount:   0,
		PairIPPort:  "",
		PairCode:    "",
		ActiveField: 0,
	}
}

func (m Model) handleWizardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "1", "u":
		m.Wizard.Mode = WizardUSB
		m.Wizard.ErrorMsg = ""
		return m, tea.Batch(LoadDevicesCmd(m.Executor), TickCmd())

	case "2", "w":
		m.Wizard.Mode = WizardWireless
		m.Wizard.ErrorMsg = ""
		return m, nil

	case "tab", "down":
		if m.Wizard.Mode == WizardWireless {
			m.Wizard.ActiveField = (m.Wizard.ActiveField + 1) % 2
		}
		return m, nil

	case "shift+tab", "up":
		if m.Wizard.Mode == WizardWireless {
			m.Wizard.ActiveField = (m.Wizard.ActiveField + 1) % 2
		}
		return m, nil

	case "backspace":
		if m.Wizard.Mode == WizardWireless {
			if m.Wizard.ActiveField == 0 && len(m.Wizard.PairIPPort) > 0 {
				m.Wizard.PairIPPort = m.Wizard.PairIPPort[:len(m.Wizard.PairIPPort)-1]
			} else if m.Wizard.ActiveField == 1 && len(m.Wizard.PairCode) > 0 {
				m.Wizard.PairCode = m.Wizard.PairCode[:len(m.Wizard.PairCode)-1]
			}
		}
		return m, nil

	case "enter":
		if m.Wizard.Mode == WizardWireless {
			if m.Wizard.PairIPPort == "" || m.Wizard.PairCode == "" {
				m.Wizard.ErrorMsg = "Both IP:Port and Pairing Code are required"
				return m, nil
			}
			m.Wizard.IsPairing = true
			m.Wizard.ErrorMsg = ""
			m.Status = "Initiating wireless pairing handshake..."
			return m, PairDeviceCmd(m.Executor, m.Wizard.PairIPPort, m.Wizard.PairCode)
		}
		return m, LoadDevicesCmd(m.Executor)

	case "r":
		m.Status = "Checking for connected devices..."
		return m, LoadDevicesCmd(m.Executor)

	default:
		if m.Wizard.Mode == WizardWireless && len(msg.Runes) > 0 {
			char := string(msg.Runes)
			if m.Wizard.ActiveField == 0 {
				if len(m.Wizard.PairIPPort) < 26 {
					m.Wizard.PairIPPort += char
				}
			} else {
				if len(m.Wizard.PairCode) < 6 {
					m.Wizard.PairCode += char
				}
			}
		}
		return m, nil
	}
}
