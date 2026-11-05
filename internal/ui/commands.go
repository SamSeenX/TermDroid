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
			MaxSize:       1200,
			BitRate:       "6M",
			MaxFps:        60,
			VideoCodec:    "h264",
			WindowTitle:   "TermDroid Desktop Stream",
		}
		_, err := exec.StartMirrorProcess(opts)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Scrcpy launch failed: %v", err))
		}
		return StatusMsg("Mirror streaming window opened (1200p, 6M, 60fps)")
	}
}

func StartRecordingCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		_, path, err := exec.StartRecordingProcess("")
		if err != nil {
			return StatusMsg(fmt.Sprintf("Recording failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Recording to: %s", path))
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

func PairDeviceCmd(exec *bridge.Executor, ipPort, code string) tea.Cmd {
	return func() tea.Msg {
		res, err := exec.PairDevice(ipPort, code)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Pairing failed: %v", err))
		}
		return PairingSuccessMsg(res)
	}
}

func PushClipboardCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		res, err := exec.PushClipboardToDevice()
		if err != nil {
			return StatusMsg(fmt.Sprintf("Clipboard push error: %v", err))
		}
		return StatusMsg(res)
	}
}

func PullClipboardCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		res, err := exec.PullClipboardFromDevice()
		if err != nil {
			return StatusMsg(fmt.Sprintf("Clipboard pull error: %v", err))
		}
		return StatusMsg(res)
	}
}

func PurgeCacheCmd(exec *bridge.Executor) tea.Cmd {
	return func() tea.Msg {
		res, err := storage.PurgeAppCache(exec)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Cache cleanup failed: %v", err))
		}
		return StatusMsg(res)
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

func ExecuteKeypadCmd(exec *bridge.Executor, action bridge.KeyActionItem) tea.Cmd {
	return func() tea.Msg {
		res, err := exec.ExecuteKeyAction(action)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Key error: %v", err))
		}
		return StatusMsg(res)
	}
}

func DisconnectDeviceCmd(exec *bridge.Executor, target string) tea.Cmd {
	return func() tea.Msg {
		err := exec.DisconnectDevice(target)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Disconnect error: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Disconnected %s", target))
	}
}
