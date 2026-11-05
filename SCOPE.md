# TermDroid Project Scope

## In Scope

### Phase 1: Core Foundation & Device Multiplexing
* ADB and Scrcpy environment bootstrapping and dependency checks.
* Device listing, state detection (USB vs TCP/IP), and automatic wireless switching.
* Process lifecycle supervisor with signal handling.

### Phase 2: Interactive TUI Dashboard
* Tabbed Bubble Tea interface (Mirror, Storage, Processes, Logcat, Settings).
* Status HUD displaying battery percentage, charging state, IP, and device name.
* Keystroke-driven remote actions (Screen off, Lock, Home, Back, Volume).

### Phase 3: Mobile Storage Analyzer (DiskMan Engine)
* Ranked directory tree scan of `/sdcard/`.
* Identification of heavy folders (WhatsApp, DCIM, Download, Telegram).
* Cache cleaner and large file discovery with interactive confirmation.

### Phase 4: Telemetry & Developer Utilities
* Live filterable Logcat viewer.
* Running process viewer with PID termination.
* Package sideloading and package management.

## Out of Scope
* Modifying closed-source OEM kernels or breaking Android security policies.
* Bypassing Android OS hardware DRM flags (`FLAG_SECURE`) for protected video feeds.
* Direct USB hardware driver development from scratch.
