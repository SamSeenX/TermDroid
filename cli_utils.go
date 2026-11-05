package main

import (
	"fmt"
	"termdroid/internal/bridge"
)

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
