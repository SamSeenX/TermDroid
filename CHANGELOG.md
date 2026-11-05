# Changelog

All notable changes to TermDroid will be documented in this file.

## [v0.6.1] - 2026-11-05

### Fixed
* Added multi-tier fallback for pulling device clipboard (`dumpsys clipboard` fallback) and user guidance when Android 13-17 restricts background clipboard reads.

## [v0.6.0] - 2026-11-05

### Added
* DiskMan Subfolder Deep Navigation allowing you to hit [Enter] on any folder to dive into its subfolder size breakdown.
* Breadcrumb path tracking and hierarchical navigation history stack.
* Instant upward navigation via [Esc] or [Backspace] returning directly to parent folders.
* Scrcpy Viewer Window Shortcut Cheatsheet embedded directly on Tab 1 for quick reference (Opt+h Home, Opt+b Back, Opt+s Apps, Opt+↑/↓ Volume, Opt+p Power).

## [v0.5.6] - 2026-11-05

### Performance & UX
* Persistent in-memory storage caching across tab transitions, keeping DiskMan scan results instantly visible when switching tabs.
* Decoupled background storage scanner allowing you to navigate other tabs while a scan completes asynchronously.
* Manual refresh available via `r` on Tab 2 when you wish to force a new storage scan.

## [v0.5.5] - 2026-11-05

### Fixed
* Fixed DiskMan Storage Scanner failing on modern Android versions (Android 11 through Android 17).
* Replaced recursive `du -d 1 /sdcard` with per-directory scanning that gracefully bypasses SELinux restricted `/sdcard/Android/data` subdirectories, preventing permission denied errors and command aborts.

## [v0.5.4] - 2026-11-05

### Added
* Interactive Text Injection Modal triggered via `t` hotkey, allowing you to paste or type any string on Mac and inject it directly into the focused cursor on Android (`adb shell input text`).

## [v0.5.3] - 2026-11-05

### Performance & Latency
* Reverted dumpsys lockscreen inspection overhead on unlock macro.
* Restored zero-latency instant unlock sequence with direct wake, swipe, and PIN injection in <200ms without background verification lag.

## [v0.5.2] - 2026-11-05

### Added
* Automatic macOS Preview launcher upon capturing screenshots via `s` or `./termdroid screenshot`, opening the saved Desktop PNG instantly for review.

## [v0.5.1] - 2026-11-05

### UX & Keybindings
* Verified global universality of `r` (refresh current view data) and `q` / `ctrl+c` (instant quit) across all 6 tabs, the USB connection wizard, and modal input overlays.

## [v0.5.0] - 2026-11-05

### Added
* Tab 6: Live Notification Stream Deck parsing active notifications, package sources, titles, and body snippets (`dumpsys notification --noredact`).
* Headless Mobile Audio Forwarder streaming phone audio to macOS speakers via low-latency Opus (`scrcpy --no-video --audio-codec=opus --audio-buffer=40`).
* Hotkey `o` on Tab 1 to trigger instant phone-to-Mac audio forwarding.
* Headless CLI command `./termdroid audio` to stream mobile sound directly from the terminal.

## [v0.4.6] - 2026-11-05

### Added
* Charging source telemetry in HUD header detecting AC Fast, USB, and Wireless charging states.
* Connected Wi-Fi network SSID detection (`cmd wifi status`) displayed prominently in the telemetry HUD.

## [v0.4.5] - 2026-11-05

### UI & Styling
* Revamped visual aesthetics with a cyber-cockpit high-contrast theme featuring electric cyan, neon violet, and crimson accents.
* Enhanced header bar with bold inverted brand badge and dynamic battery charging and temperature telemetry glyphs.
* Styled hotkey badges in tab views with distinct contrast tags.
* Upgraded Hardware Keypad table cards with double borders and cursor pointers (▶).

## [v0.4.4] - 2026-11-05

### Security & Safety
* Added Guarded Unlock Protection querying `dumpsys power` and `dumpsys window` before injecting credentials.
* Suppressed PIN injection when the device is already awake and unlocked, preventing accidental credential leakage into open search bars or text fields.

## [v0.4.3] - 2026-11-05

### Added
* 1-Tap Lockscreen PIN Bypass Macro binding `l` to wake screen, swipe up, inject saved PIN, and submit in <200ms.
* Secure local PIN persistence in `~/.config/termdroid/config.json` with 0600 file permissions.
* Interactive PIN configuration modal prompt with masked password input (press `L` to edit or wipe).
* Headless CLI unlock command (`termdroid unlock [PIN]`).
* Mute Volume hotkey `u` preserved for audio control while `l` handles instant unlock.

### Refactored
* Modularized UI directional navigation and tab transitions into `internal/ui/navigation.go`.
* Extracted PIN modal prompt into `internal/ui/pin_modal.go` ensuring all files strictly remain under 180 lines.

## [v0.4.2] - 2026-11-05

### Fixed
* Resolved frozen video stream issue over Wi-Fi on high-resolution devices like Pixel 6 Pro.
* Capped wireless stream bounds to 1200p resolution, 6Mbps video bitrate, and 60 FPS.
* Enforced explicit H.264 video codec encoding preventing packet loss decoder stalls.

## [v0.4.1] - 2026-11-05

### Added
* Headless coordinate tap simulation (`termdroid tap <X> <Y>`).
* Headless gesture swipe simulation with custom duration (`termdroid swipe <X1> <Y1> <X2> <Y2> [duration]`).

## [v0.4.0] - 2026-11-05

### Added
* Real-time MP4 screen recording hotkey (v) saving directly to Desktop.
* Bi-directional clipboard synchronization (c to push Mac clipboard, y to pull phone clipboard).
* DiskMan application cache purger hotkey (c on Tab 2) removing junk app caches.
* CLI APK sideloading command (`termdroid install <path.apk>`).
* CLI URL dispatcher command (`termdroid open <url>`).
* Reverse port forwarding command (`termdroid reverse <port>`).

### Refactored
* Modularized CLI subcommand dispatching into `cli_commands.go` and `cli_utils.go` maintaining strict compliance with the 200-line rule.

## [v0.3.1] - 2026-11-05

### Added
* Compact 3-Column side-by-side grid layout for Tab 5 fitting all 33 simulation buttons on one page without scrolling.
* Discrete Call Volume (+/-) stream controls targeting Stream 0 with on-screen visual slider feedback.
* Discrete Media Volume (+/-) and Ringtone (+/-) stream controls targeting Streams 3 and 2.
* Quick 1-touch hotkeys for frequent actions (h, b, e, o, +, -, c, v, u, Space, n, p, j, g, s) while using arrow keys for full grid navigation.

### Refactored
* Extracted keypad action dispatcher into `internal/ui/keypad_actions.go` keeping all files under 150 lines.

## [v0.3.0] - 2026-11-05

### Added
* Dedicated Hardware Keypad & Control Deck on Tab 5 featuring 27 categorized button simulations.
* Support for Navigation, Audio & Volume, Media Playback, Status Panels, D-Pad Remote, and Text Clipboard commands.
* Multi-device switcher hotkey (d) to cycle active target devices.
* Wireless device forget/disconnect hotkey (x) triggering adb disconnect.
* On-demand Add Device modal trigger (a) to summon the Connection Wizard anytime.

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
