package ui

import (
	"fmt"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab == TabActions {
		return m.handleKeypadEnter()
	}

	if m.ActiveTab == TabStorage && m.Storage != nil && len(m.Storage.Items) > 0 {
		if m.StorageCursor < len(m.Storage.Items) {
			selected := m.Storage.Items[m.StorageCursor]
			m.StorageStack = append(m.StorageStack, m.StoragePath)
			m.StoragePath = selected.Path
			m.Storage = nil
			m.StorageCursor = 0
			m.IsScanningStorage = true
			m.IsLoading = true
			m.Status = fmt.Sprintf("Entering %s...", selected.Path)
			return m, LoadStorageCmd(m.Executor, selected.Path)
		}
	}
	return m, nil
}

func (m Model) handleBackNavKey() (tea.Model, tea.Cmd) {
	if m.ConfirmDelete {
		m.ConfirmDelete = false
		m.DeleteTarget = ""
		m.Status = "Deletion cancelled"
		return m, nil
	}

	if m.ActiveTab == TabStorage && len(m.StorageStack) > 0 {
		prev := m.StorageStack[len(m.StorageStack)-1]
		m.StorageStack = m.StorageStack[:len(m.StorageStack)-1]
		m.StoragePath = prev
		m.Storage = nil
		m.StorageCursor = 0
		m.IsScanningStorage = true
		m.IsLoading = true
		m.Status = fmt.Sprintf("Returning to %s...", prev)
		return m, LoadStorageCmd(m.Executor, prev)
	}
	return m, nil
}

func (m Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if m.ActiveTab != TabStorage || m.Storage == nil || len(m.Storage.Items) == 0 {
		return m, nil
	}

	targets := m.getMarkedOrCursorPaths()
	if len(targets) == 0 {
		return m, nil
	}

	m.ConfirmDelete = true
	if len(targets) == 1 {
		m.DeleteTarget = targets[0]
		m.Status = fmt.Sprintf("Delete %s? Press [y] to confirm, [Esc] to cancel", targets[0])
	} else {
		m.DeleteTarget = fmt.Sprintf("%d items", len(targets))
		m.Status = fmt.Sprintf("Delete %d items? Press [y] to confirm, [Esc] to cancel", len(targets))
	}
	return m, nil
}

func (m Model) handleConfirmDelete() (tea.Model, tea.Cmd) {
	if !m.ConfirmDelete {
		return m, nil
	}

	targets := m.getMarkedOrCursorPaths()
	m.ConfirmDelete = false
	m.DeleteTarget = ""
	m.SelectedFiles = make(map[string]bool)

	if len(targets) == 0 {
		return m, nil
	}

	targetMap := make(map[string]bool)
	for _, t := range targets {
		targetMap[t] = true
	}

	var remaining []storage.StorageItem
	for _, it := range m.Storage.Items {
		if !targetMap[it.Path] {
			remaining = append(remaining, it)
		}
	}
	m.Storage.Items = remaining

	if m.StorageCursor >= len(remaining) && len(remaining) > 0 {
		m.StorageCursor = len(remaining) - 1
	}

	m.Status = fmt.Sprintf("Deleting %d item(s) from device...", len(targets))
	return m, BatchDeleteStorageCmd(m.Executor, targets)
}

func (m Model) handleActionShortcut(keyStr string) (tea.Model, tea.Cmd) {
	if act, ok := findActionByShortcut(keyStr); ok {
		m.Status = fmt.Sprintf("Shortcut fired: %s...", act.Name)
		return m, ExecuteKeypadCmd(m.Executor, act)
	}
	return m, nil
}

func handleMirrorKeys(m Model, keyStr string) (Model, tea.Cmd, bool) {
	if m.ActiveTab != TabMirror {
		return m, nil, false
	}
	switch keyStr {
	case "m":
		m.Status = "Launching Scrcpy live mirror..."
		return m, LaunchMirrorCmd(m.Executor, false), true
	case "p":
		m.Status = "Launching Scrcpy in Presentation Mode..."
		return m, LaunchMirrorCmd(m.Executor, true), true
	case "v":
		m.Status = "Initiating live screen recording..."
		return m, StartRecordingCmd(m.Executor), true
	case "o":
		m.Status = "Launching mobile audio forwarder..."
		return m, LaunchAudioStreamCmd(m.Executor), true
	case "c":
		m.Status = "Pushing host clipboard to phone..."
		return m, PushClipboardCmd(m.Executor), true
	case "s":
		m.Status = "Capturing screen frame..."
		return m, TakeScreenshotCmd(m.Executor), true
	case "k":
		m.Status = "Injecting wake and swipe unlock gesture..."
		return m, UnlockDeviceCmd(m.Executor, ""), true
	default:
		return m, nil, false
	}
}
