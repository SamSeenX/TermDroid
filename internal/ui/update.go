package ui

import (
	"fmt"
	"termdroid/internal/bridge"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case StatusMsg:
		m.Status = string(msg)
		m.IsLoading = false
		m.Wizard.IsPairing = false
		return m, nil

	case PairingSuccessMsg:
		m.Status = string(msg)
		m.Wizard.IsPairing = false
		return m, LoadDevicesCmd(m.Executor)

	case TickMsg:
		m.Wizard.TickCount++
		if len(m.Devices) == 0 {
			return m, tea.Batch(LoadDevicesCmd(m.Executor), TickCmd())
		}
		return m, TickCmd()

	case DevicesLoadedMsg:
		return m.handleDevicesLoaded(msg)

	case VitalsLoadedMsg:
		m.Vitals = msg
		return m, nil

	case StorageLoadedMsg:
		m.Storage = msg
		m.IsLoading = false
		m.Status = fmt.Sprintf("Storage: %s used of %s", m.Storage.UsedSpace, m.Storage.TotalSpace)
		return m, nil

	case ProcessesLoadedMsg:
		m.Processes = msg
		m.IsLoading = false
		m.Status = fmt.Sprintf("Loaded %d active processes", len(m.Processes))
		return m, nil

	case LogcatLoadedMsg:
		m.LogLines = msg
		m.IsLoading = false
		m.Status = fmt.Sprintf("Retrieved %d recent log lines", len(m.LogLines))
		return m, nil

	case NotificationsLoadedMsg:
		m.Notifications = msg
		m.IsLoading = false
		m.Status = fmt.Sprintf("Retrieved %d active notifications", len(m.Notifications))
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m Model) handleDevicesLoaded(msg DevicesLoadedMsg) (tea.Model, tea.Cmd) {
	m.IsLoading = false
	var authorized []bridge.DeviceInfo
	var unauthorized []bridge.DeviceInfo

	for _, d := range msg {
		if d.State == "device" {
			authorized = append(authorized, d)
		} else if d.State == "unauthorized" {
			unauthorized = append(unauthorized, d)
		}
	}

	if len(authorized) > 0 {
		m.Devices = authorized
		m.ActiveDevice = 0
		m.Executor.DeviceID = authorized[0].ID
		mode := "USB"
		if authorized[0].IsWireless {
			mode = "Wi-Fi"
		}
		m.Status = fmt.Sprintf("Connected to %s via %s", authorized[0].Model, mode)
		return m, LoadVitalsCmd(m.Executor)
	}

	m.Devices = nil
	if len(unauthorized) > 0 {
		dev := unauthorized[0]
		m.Status = fmt.Sprintf("Found %s: UNAUTHORIZED. Tap 'Allow' on phone screen.", dev.Model)
		m.Wizard.ErrorMsg = "Device unauthorized. Check phone and tap 'Allow USB debugging'."
		return m, nil
	}

	if m.Wizard.Mode == WizardUSB {
		m.Status = "Monitoring for USB cable connection..."
	}
	return m, nil
}
