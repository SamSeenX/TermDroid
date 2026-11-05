package main

import (
	"flag"
	"fmt"
	"termdroid/internal/bridge"
	"termdroid/internal/storage"
)

func handleCLICommand(command string, args []string) {
	exec := bridge.NewExecutor("")

	switch command {
	case "devices":
		handleDevicesCmd(exec)

	case "mirror":
		handleMirrorCmd(exec, args)

	case "record":
		handleRecordCmd(exec, args)

	case "scan":
		handleScanCmd(exec)

	case "shot", "screenshot":
		handleScreenshotCmd(exec, args)

	case "clipboard":
		handleClipboardCmd(exec, args)

	case "install":
		handleInstallCmd(exec, args)

	case "open":
		handleOpenCmd(exec, args)

	case "reverse":
		handleReverseCmd(exec, args)

	case "wifi":
		handleWifiCmd(exec)

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
	}
}

func handleDevicesCmd(exec *bridge.Executor) {
	devs, err := exec.ListDevices()
	if err != nil {
		fmt.Printf("Error querying devices: %v\n", err)
		return
	}
	fmt.Printf("Connected Android Devices (%d):\n", len(devs))
	for _, d := range devs {
		mode := "USB"
		if d.IsWireless {
			mode = "Wi-Fi"
		}
		notice := ""
		if d.State == "unauthorized" {
			notice = " [UNAUTHORIZED: Tap Allow on phone]"
		}
		fmt.Printf("- %s [%s] (%s)%s\n", d.Model, d.ID, mode, notice)
	}
}

func handleMirrorCmd(exec *bridge.Executor, args []string) {
	fs := flag.NewFlagSet("mirror", flag.ExitOnError)
	screenOff := fs.Bool("screen-off", false, "Turn off device screen")
	recordPath := fs.String("record", "", "Path to save MP4 recording")
	fs.Parse(args)

	opts := bridge.MirrorOptions{
		TurnScreenOff: *screenOff,
		StayAwake:     true,
		RecordPath:    *recordPath,
	}
	fmt.Println("Launching Scrcpy mirror stream...")
	cmd, err := exec.StartMirrorProcess(opts)
	if err != nil {
		fmt.Printf("Error starting mirror: %v\n", err)
		return
	}
	cmd.Wait()
}

func handleRecordCmd(exec *bridge.Executor, args []string) {
	dest := ""
	if len(args) > 0 {
		dest = args[0]
	}
	cmd, path, err := exec.StartRecordingProcess(dest)
	if err != nil {
		fmt.Printf("Recording launch failed: %v\n", err)
		return
	}
	fmt.Printf("Recording live screen stream to: %s\n(Close window to finalize video)\n", path)
	cmd.Wait()
}

func handleScanCmd(exec *bridge.Executor) {
	fmt.Println("Scanning Android internal storage (/sdcard)...")
	summary, err := storage.ScanStorageTree(exec, "/sdcard")
	if err != nil {
		fmt.Printf("Storage scan error: %v\n", err)
		return
	}
	fmt.Printf("Storage: %s used of %s total (%s)\n\n", summary.UsedSpace, summary.TotalSpace, summary.Percent)
	for _, item := range summary.Items {
		fmt.Printf("%-10s  %-20s  %s\n", item.SizeHuman, item.Category, item.Path)
	}
}

func handleScreenshotCmd(exec *bridge.Executor, args []string) {
	dest := ""
	if len(args) > 0 {
		dest = args[0]
	}
	saved, err := exec.TakeScreenshot(dest)
	if err != nil {
		fmt.Printf("Screenshot error: %v\n", err)
		return
	}
	fmt.Printf("Screenshot saved to: %s\n", saved)
}
