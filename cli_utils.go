package main

import (
	"fmt"
	"strconv"
	"termdroid/internal/bridge"
)

func handleTapCmd(exec *bridge.Executor, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: termdroid tap <X> <Y>")
		return
	}
	x, err1 := strconv.Atoi(args[0])
	y, err2 := strconv.Atoi(args[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Error: coordinates must be integers")
		return
	}
	if err := exec.Tap(x, y); err != nil {
		fmt.Printf("Tap injection error: %v\n", err)
		return
	}
	fmt.Printf("Injected tap at (%d, %d)\n", x, y)
}

func handleSwipeCmd(exec *bridge.Executor, args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: termdroid swipe <X1> <Y1> <X2> <Y2> [duration_ms]")
		return
	}
	x1, _ := strconv.Atoi(args[0])
	y1, _ := strconv.Atoi(args[1])
	x2, _ := strconv.Atoi(args[2])
	y2, _ := strconv.Atoi(args[3])
	duration := 300
	if len(args) >= 5 {
		duration, _ = strconv.Atoi(args[4])
	}
	if err := exec.Swipe(x1, y1, x2, y2, duration); err != nil {
		fmt.Printf("Swipe injection error: %v\n", err)
		return
	}
	fmt.Printf("Injected swipe from (%d, %d) to (%d, %d) in %dms\n", x1, y1, x2, y2, duration)
}

func handleAudioCmd(exec *bridge.Executor) {
	fmt.Println("Launching headless mobile audio stream to Mac speakers (Opus, low-latency)...")
	cmd, err := exec.StartAudioStreamProcess()
	if err != nil {
		fmt.Printf("Audio launch error: %v\n", err)
		return
	}
	fmt.Println("Audio forwarder active. Press Ctrl+C to disconnect audio stream.")
	cmd.Wait()
}

func handleUnlockCmd(exec *bridge.Executor, args []string) {
	pin := ""
	if len(args) > 0 {
		pin = args[0]
		bridge.SaveConfig(bridge.Config{SavedPIN: pin})
	} else {
		cfg := bridge.LoadConfig()
		pin = cfg.SavedPIN
	}

	if pin != "" {
		fmt.Println("Waking device, swiping up, and injecting PIN...")
	} else {
		fmt.Println("No PIN specified. Performing wake and swipe up gesture...")
	}

	if err := exec.UnlockDevice(pin); err != nil {
		fmt.Printf("Unlock error: %v\n", err)
		return
	}
	fmt.Println("Unlock sequence executed successfully.")
}

func handleClipboardCmd(exec *bridge.Executor, args []string) {
	sub := "push"
	if len(args) > 0 {
		sub = args[0]
	}
	if sub == "pull" {
		res, err := exec.PullClipboardFromDevice()
		if err != nil {
			fmt.Printf("Clipboard pull error: %v\n", err)
			return
		}
		fmt.Println(res)
		return
	}

	res, err := exec.PushClipboardToDevice()
	if err != nil {
		fmt.Printf("Clipboard push error: %v\n", err)
		return
	}
	fmt.Println(res)
}

func handleInstallCmd(exec *bridge.Executor, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: termdroid install <path/to/app.apk>")
		return
	}
	fmt.Printf("Installing %s onto device...\n", args[0])
	if err := exec.InstallAPK(args[0]); err != nil {
		fmt.Printf("Installation error: %v\n", err)
		return
	}
	fmt.Println("APK installed successfully.")
}

func handleOpenCmd(exec *bridge.Executor, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: termdroid open <https://example.com>")
		return
	}
	if err := exec.OpenURL(args[0]); err != nil {
		fmt.Printf("Failed to open URL: %v\n", err)
		return
	}
	fmt.Printf("Dispatched URL %s to mobile browser.\n", args[0])
}

func handleReverseCmd(exec *bridge.Executor, args []string) {
	port := "8080"
	if len(args) > 0 {
		port = args[0]
	}
	if err := exec.ReversePort(port); err != nil {
		fmt.Printf("Reverse port forwarding error: %v\n", err)
		return
	}
	fmt.Printf("Reverse port forwarded: localhost:%s on phone reaches your Mac\n", port)
}

func handleWifiCmd(exec *bridge.Executor) {
	fmt.Println("Configuring wireless TCP/IP bridge...")
	addr, err := exec.SwitchToWireless("5555")
	if err != nil {
		fmt.Printf("Wireless switch failed: %v\n", err)
		return
	}
	fmt.Printf("Successfully connected wirelessly to %s\n", addr)
}
