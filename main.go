package main

import (
	"flag"
	"fmt"
	"os"
	"termdroid/internal/bridge"
	"termdroid/internal/storage"
	"termdroid/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) > 1 {
		handleCLICommand(os.Args[1], os.Args[2:])
		return
	}

	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing TermDroid: %v\n", err)
		os.Exit(1)
	}
}

func handleCLICommand(command string, args []string) {
	exec := bridge.NewExecutor("")

	switch command {
	case "devices":
		devs, err := exec.ListDevices()
		if err != nil {
			fmt.Printf("Error querying devices: %v\n", err)
			return
		}
		fmt.Printf("Connected Android Devices (%d):\n", len(devs))
		for _, d := range devs {
			mode := "USB"
			if d.IsWireless {
				mode = "Wi-Fi TCP/IP"
			}
			fmt.Printf("- %s [%s] (%s)\n", d.Model, d.ID, mode)
		}

	case "mirror":
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

	case "scan":
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

	case "shot", "screenshot":
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

	case "wifi":
		fmt.Println("Configuring wireless TCP/IP bridge...")
		addr, err := exec.SwitchToWireless("5555")
		if err != nil {
			fmt.Printf("Wireless switch failed: %v\n", err)
			return
		}
		fmt.Printf("Successfully connected wirelessly to %s\n", addr)

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("TermDroid - macOS Android Control Cockpit & Storage Analyzer")
	fmt.Println("\nUsage:")
	fmt.Println("  termdroid                 Launch interactive Bubble Tea TUI dashboard")
	fmt.Println("  termdroid devices         List all attached USB and Wi-Fi devices")
	fmt.Println("  termdroid mirror          Start low-latency Scrcpy mirror window")
	fmt.Println("  termdroid scan            Run DiskMan storage inspection on device")
	fmt.Println("  termdroid screenshot      Capture current screen frame to desktop")
	fmt.Println("  termdroid wifi            Switch connected USB device to wireless TCP/IP")
}
