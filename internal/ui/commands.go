package ui

import (
	"fmt"
	"termdroid/internal/bridge"
	"termdroid/internal/diagnostics"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func LoadDevicesCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		devices, err := exec.ListDevices()
		if err != nil {
			return StatusMsg(fmt.Sprintf("Error listing devices: %v", err))
		}
		return DevicesLoadedMsg(devices)
	}
}

func LoadVitalsCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		v, err := diagnostics.FetchVitals(exec)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Vitals error: %v", err))
		}
		return VitalsLoadedMsg(v)
	}
}

func LoadStorageCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		s, err := storage.ScanStorageTree(exec, "/sdcard")
		if err != nil {
			return StatusMsg(fmt.Sprintf("Storage scan error: %v", err))
		}
		return StorageLoadedMsg(s)
	}
}

func LoadProcessesCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		procs, err := diagnostics.ListProcesses(exec)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Process list error: %v", err))
		}
		return ProcessesLoadedMsg(procs)
	}
}

func LoadLogcatCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		logs, err := diagnostics.FetchRecentLogs(exec, 40, diagnostics.LogInfo)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Logcat error: %v", err))
		}
		return LogcatLoadedMsg(logs)
	}
}

func LaunchMirrorCmd(exec *bridge.Executor, turnScreenOff bool) tea.Cmd {
	return func() tea.Msg {
		opts := bridge.MirrorOptions{
			TurnScreenOff: turnScreenOff,
			StayAwake:     true,
			AlwaysOnTop:   true,
			WindowTitle:   "TermDroid Desktop Stream",
		}
		_, err := exec.StartMirrorProcess(opts)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Scrcpy launch failed: %v", err))
		}
		return StatusMsg("Mirror streaming window opened")
	}
}

func SwitchWifiCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		target, err := exec.SwitchToWireless("5555")
		if err != nil {
			return StatusMsg(fmt.Sprintf("Wireless switch failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Connected wirelessly to %s", target))
	}
}

func TakeScreenshotCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		path, err := exec.TakeScreenshot("")
		if err != nil {
			return StatusMsg(fmt.Sprintf("Screenshot failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Screenshot saved: %s", path))
	}
}

func SendKeyCmd(exec *bridge.Executor, key bridge.KeyAction) tea.Cmd {
	return func() tea.Msg {
		err := exec.SendKeyEvent(key)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Keyevent error: %v", err))
		}
		return StatusMsg("Key command sent")
	}
}
