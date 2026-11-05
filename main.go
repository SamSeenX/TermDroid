package main

import (
	"fmt"
	"os"
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

func printHelp() {
	fmt.Println("TermDroid - macOS Android Control Cockpit & Storage Analyzer")
	fmt.Println("\nUsage:")
	fmt.Println("  termdroid                 Launch interactive Bubble Tea TUI dashboard")
	fmt.Println("  termdroid devices         List all attached USB and Wi-Fi devices")
	fmt.Println("  termdroid mirror          Start low-latency Scrcpy mirror window")
	fmt.Println("  termdroid record          Record live mirror stream to Desktop MP4")
	fmt.Println("  termdroid scan            Run DiskMan storage inspection on device")
	fmt.Println("  termdroid screenshot      Capture current screen frame to desktop")
	fmt.Println("  termdroid tap <x> <y>     Inject headless touch tap at coordinates")
	fmt.Println("  termdroid swipe x1 y1 x2 y2 Inject headless touch swipe gesture")
	fmt.Println("  termdroid unlock [PIN]    1-tap unlock macro (wake, swipe, inject PIN)")
	fmt.Println("  termdroid clipboard       Push or pull clipboard (push/pull)")
	fmt.Println("  termdroid install <apk>   Sideload APK package onto device")
	fmt.Println("  termdroid open <url>      Open web URL in default mobile browser")
	fmt.Println("  termdroid reverse <port>  Reverse port forward local Mac server to phone")
	fmt.Println("  termdroid wifi            Switch connected USB device to wireless TCP/IP")
}
