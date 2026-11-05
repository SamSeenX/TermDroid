# TermDroid Architectural Brain

## 1. Architectural Philosophy
TermDroid is an independent macOS terminal dashboard and CLI suite for Android device management. It rejects the design of being merely a thin Scrcpy wrapper. Instead, TermDroid delivers deep hardware control, Linux storage analytics, process monitoring, and telemetry through native ADB pipelines, while maintaining Scrcpy as a modular subsystem for low-latency visual mirroring.

## 2. Core Decisions
* Language: Go 1.26+ for single-binary distribution, zero runtime dependencies, and high-concurrency goroutine management.
* TUI Framework: Charmbracelet Bubble Tea with Lip Gloss styling, Bubbles components, and Elm architecture.
* Subprocess Supervisor: Strict OS signal trapping and context cancellation to prevent orphaned Scrcpy or Logcat background tasks.
* Storage Engine: Mobile adaptation of the DiskMan engine, using bounded concurrent ADB worker pools to scan, rank, and clean mobile storage.

## 3. State Management
* Root Application State: Managed by the Bubble Tea Elm loop (`Model`, `Update`, `View`).
* Sub-View States: Partitioned into Device HUD, Storage Tree, Process Inspector, Logcat Streamer, and Mirror Controller.
* Concurrency: Asynchronous tasks run in Go commands returning typed Bubble Tea messages. Long-running streams (like Logcat) push into buffered Go channels.

## 4. Hardware and Engine Interfaces
* ADB Server: Discovered dynamically via `/opt/homebrew/bin/adb` or `$PATH`. Reuses existing daemons without disrupting running developer tools.
* Scrcpy Engine: Launched in Presentation Mode by default to preserve device battery during live mirroring sessions.
