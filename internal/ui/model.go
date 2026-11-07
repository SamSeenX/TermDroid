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
	TabNotifications
)

var TabNames = []string{
	"1. Display & Mirror",
	"2. DiskMan Storage",
	"3. Processes",
	"4. Logcat",
	"5. Hardware Keypad & Control",
	"6. Notifications Stream",
}

type DevicesLoadedMsg []bridge.DeviceInfo
type VitalsLoadedMsg *diagnostics.DeviceVitals
type StorageLoadedMsg *storage.StorageSummary
type ProcessesLoadedMsg []diagnostics.ProcessItem
type LogcatLoadedMsg []string
type NotificationsLoadedMsg []diagnostics.NotificationItem
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
	Notifications []diagnostics.NotificationItem
	StorageCursor int
	StoragePath   string
	StorageStack  []string
	ProcessCursor int
	NotifCursor   int
	KeypadCol     int
	KeypadRow     int
	Status        string
	Width         int
	Height        int
	IsLoading     bool
	Wizard        WizardState
	PinModalOpen  bool
	PinInput      string
	SavedPIN      string
	TextModalOpen bool
	TextInput     string
	IsScanningStorage bool
	DeleteTarget  string
	ConfirmDelete bool
	SelectedFiles map[string]bool
	ClipboardPaths []string
	ClipboardIsCut bool
	StorageModalOpen bool
	StorageModalType StorageModalAction
	StorageModalInput string
}

func InitialModel() Model {
	exec := bridge.NewExecutor("")
	cfg := bridge.LoadConfig()
	return Model{
		Executor:      exec,
		ActiveTab:     TabMirror,
		Status:        "Checking for connected Android devices...",
		StorageCursor: 0,
		ProcessCursor: 0,
		KeypadCol:     0,
		KeypadRow:     0,
		IsLoading:     false,
		Wizard:        InitialWizardState(),
		SavedPIN:      cfg.SavedPIN,
		SelectedFiles: make(map[string]bool),
		StorageModalType: ModalNone,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		LoadDevicesCmd(m.Executor),
		TickCmd(),
		tea.EnterAltScreen,
	)
}
