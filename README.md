# TermDroid

TermDroid is a high-performance terminal cockpit for macOS, Linux, and Windows that connects, mirrors, monitors, and manages Android devices over USB and Wi-Fi.

## Key Features

* Low-Latency Screen Mirroring: Real-time interactive desktop window powered by Scrcpy with optional headless or screen-off presentation modes.
* Mobile Storage Inspector: Ported from the DiskMan Mac storage engine. Concurrently scans `/sdcard/`, discovers bulky media, identifies hidden caches, and facilitates rapid disk cleanup.
* Device Telemetry HUD: Real-time tracking of battery vitals, charging rates, CPU load, and memory usage.
* Process Monitor & Control: List and filter running processes on Android with one-keystroke termination of runaway background apps.
* Filterable Logcat Stream: Real-time, regex-searchable system log streamer built directly into a dedicated terminal viewport.
* Smart Connectivity: Instant switching between wired USB and wireless TCP/IP mode with connection discovery.

## Prerequisites

* macOS, Linux (with Wayland or X11), or Windows 10/11
* Android device with Developer Options and USB Debugging enabled
* External tools:
  * ADB (Android Platform Tools)
  * Scrcpy

## Installation

### Option 1: Download Pre-compiled Binary (Recommended)

Download the ready-to-run executable for your operating system from [Releases](https://github.com/SamSeenX/TermDroid/releases):

* **macOS Apple Silicon**: `termdroid-v*-darwin-arm64.tar.gz`
* **macOS Intel**: `termdroid-v*-darwin-amd64.tar.gz`
* **Linux 64-bit**: `termdroid-v*-linux-amd64.tar.gz`
* **Linux ARM64**: `termdroid-v*-linux-arm64.tar.gz`
* **Windows 64-bit**: `termdroid-v*-windows-amd64.zip`

Extract the archive and run the executable directly.

### Option 2: Build from Source

```bash
# Build the binary
go build -o termdroid .
```

## Quick Start

# Launch interactive TUI
./termdroid

# Direct CLI commands
./termdroid mirror
./termdroid scan
./termdroid wifi
```

## Navigation Keys

* `Tab` / `Shift+Tab`: Switch between view tabs
* `j` / `k` or `Up` / `Down`: Move selection cursor
* `Enter`: Execute action or drill into folder
* `Esc`: Go back or cancel operation
* `q` or `Ctrl+C`: Exit TermDroid
