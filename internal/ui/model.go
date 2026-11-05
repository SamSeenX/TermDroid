package ui

import (
	"termdroid/internal/bridge"
	"termdroid/internal/diagnostics"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

const (
	TabMirror Tab = iota
	TabStorage
	TabProcesses
	TabLogcat
	TabActions
)

var TabNames = []string{
	"1. Display & Mirror",
	"2. DiskMan Storage",
	"3. Processes",
	"4. Logcat",
	"5. Hardware & Quick Actions",
}

type DevicesLoadedMsg []bridge.DeviceInfo
type VitalsLoadedMsg *diagnostics.DeviceVitals
type StorageLoadedMsg *storage.StorageSummary
type ProcessesLoadedMsg []diagnostics.ProcessItem
type LogcatLoadedMsg []string
type StatusMsg string
type PairingSuccessMsg string

type Model struct {
	Executor      *bridge.Executor
	Devices       []bridge.DeviceInfo
	ActiveDevice  int
	ActiveTab     Tab
	Vitals        *diagnostics.DeviceVitals
	Storage       *storage.StorageSummary
	Processes     []diagnostics.ProcessItem
	LogLines      []string
	StorageCursor int
	ProcessCursor int
	Status        string
	Width         int
	Height        int
	IsLoading     bool
	Wizard        WizardState
}

func InitialModel() Model {
	exec := bridge.NewExecutor("")
	return Model{
		Executor:      exec,
		ActiveTab:     TabMirror,
		Status:        "Checking for connected Android devices...",
		StorageCursor: 0,
		ProcessCursor: 0,
		IsLoading:     false,
		Wizard:        InitialWizardState(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		LoadDevicesCmd(m.Executor),
		TickCmd(),
		tea.EnterAltScreen,
	)
}
