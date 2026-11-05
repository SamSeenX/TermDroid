# TermDroid

TermDroid is a high-performance terminal cockpit for macOS that connects, mirrors, monitors, and manages Android devices over USB and Wi-Fi.

## Key Features

* Low-Latency Screen Mirroring: Real-time interactive desktop window powered by Scrcpy with optional headless or screen-off presentation modes.
* Mobile Storage Inspector: Ported from the DiskMan Mac storage engine. Concurrently scans `/sdcard/`, discovers bulky media, identifies hidden caches, and facilitates rapid disk cleanup.
* Device Telemetry HUD: Real-time tracking of battery vitals, charging rates, CPU load, and memory usage.
* Process Monitor & Control: List and filter running processes on Android with one-keystroke termination of runaway background apps.
* Filterable Logcat Stream: Real-time, regex-searchable system log streamer built directly into a dedicated terminal viewport.
* Smart Connectivity: Instant switching between wired USB and wireless TCP/IP mode with connection discovery.

## Prerequisites

* macOS 12+ (Apple Silicon or Intel)
* Android device with Developer Options and USB Debugging enabled
* External tools:
  * ADB (Android Platform Tools): `brew install android-platform-tools`
  * Scrcpy: `brew install scrcpy`

## Quick Start

```bash
# Build the binary
go build -o termdroid main.go

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
