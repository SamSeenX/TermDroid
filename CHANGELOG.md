# Changelog

All notable changes to TermDroid will be documented in this file.

## [v0.2.2] - 2026-11-05

### Added
* Automatic transition directly to dashboard if an authorized device is already connected on startup.
* Real-time authorization monitor detecting devices in unauthorized state with on-screen prompt to tap Allow.
* Chained automatic device discovery immediately upon successful wireless pairing.
* Support for mDNS wireless connection identifiers (`._adb-tls-connect._tcp`).

### Refactored
* Extracted keypress management into `internal/ui/keys.go` adhering to the proactive 180-line modular rule.

## [v0.2.1] - 2026-11-05

### Fixed
* Isolated keystroke handling between USB Watcher and Wireless text input modes.
* Prevented digits (such as 1 or 2) and letters from triggering USB mode switches during IP address entry.
* Assigned Esc key for explicit, unambiguous return to USB mode from Wireless pairing.

## [v0.2.0] - 2026-11-05

### Added
* Zero-device Connection Wizard with dual setup modes (USB and Wireless).
* Background USB device watcher with animated polling ticker and checklist.
* Android 11+ Wireless Debugging form with pairing code support.
* ADB pairing command execution pipeline.

### Fixed
* Corrected Lip Gloss rendering alignment artifacts across all tab views.
* Enhanced case-insensitive directory category matching in DiskMan storage engine.

## [v0.1.0] - 2026-11-05

### Added
* Project initialization with standard repository architecture.
* Elephant-Goldfish feasibility study and architectural roadmap.
* Subprocess supervision engine for ADB and Scrcpy with signal trapping.
* Bubble Tea Elm architecture TUI skeleton with multi-tab layout.
* Port of DiskMan storage inspection engine for Android internal filesystems.
* Device telemetry, process inspection, and wireless switching modules.
