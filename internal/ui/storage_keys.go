package ui

import (
	"fmt"
	"path/filepath"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func handleStorageKeys(m Model, keyStr string) (Model, tea.Cmd, bool) {
	if m.ActiveTab != TabStorage {
		return m, nil, false
	}

	switch keyStr {
	case " ":
		return m.handleToggleMark()

	case "a":
		return m.handleSelectAll()

	case "d", "P":
		return m.handlePullSelected()

	case "u":
		m.StorageModalOpen = true
		m.StorageModalType = ModalPush
		m.StorageModalInput = ""
		m.Status = "Enter computer file path to upload..."
		return m, nil, true

	case "c":
		return m.handleCopy()

	case "m":
		return m.handleCut()

	case "p":
		return m.handlePaste()

	case "n":
		m.StorageModalOpen = true
		m.StorageModalType = ModalMkdir
		m.StorageModalInput = ""
		m.Status = "Enter new directory name..."
		return m, nil, true

	case "R":
		if m.Storage == nil || len(m.Storage.Items) == 0 || m.StorageCursor >= len(m.Storage.Items) {
			return m, nil, true
		}
		currName := filepath.Base(m.Storage.Items[m.StorageCursor].Path)
		m.StorageModalOpen = true
		m.StorageModalType = ModalRename
		m.StorageModalInput = currName
		m.Status = fmt.Sprintf("Renaming %s...", currName)
		return m, nil, true

	default:
		return m, nil, false
	}
}

func (m Model) handleToggleMark() (Model, tea.Cmd, bool) {
	if m.Storage == nil || len(m.Storage.Items) == 0 || m.StorageCursor >= len(m.Storage.Items) {
		return m, nil, true
	}
	if m.SelectedFiles == nil {
		m.SelectedFiles = make(map[string]bool)
	}
	item := m.Storage.Items[m.StorageCursor]
	m.SelectedFiles[item.Path] = !m.SelectedFiles[item.Path]
	if !m.SelectedFiles[item.Path] {
		delete(m.SelectedFiles, item.Path)
	}
	m.Status = fmt.Sprintf("Selected %d item(s)", len(m.SelectedFiles))
	return m, nil, true
}

func (m Model) handleSelectAll() (Model, tea.Cmd, bool) {
	if m.Storage == nil || len(m.Storage.Items) == 0 {
		return m, nil, true
	}
	if m.SelectedFiles == nil {
		m.SelectedFiles = make(map[string]bool)
	}
	if len(m.SelectedFiles) >= len(m.Storage.Items) {
		m.SelectedFiles = make(map[string]bool)
		m.Status = "Deselected all items"
		return m, nil, true
	}
	for _, it := range m.Storage.Items {
		m.SelectedFiles[it.Path] = true
	}
	m.Status = fmt.Sprintf("Selected all %d items", len(m.Storage.Items))
	return m, nil, true
}

func (m Model) getMarkedOrCursorPaths() []string {
	var paths []string
	for p, marked := range m.SelectedFiles {
		if marked {
			paths = append(paths, p)
		}
	}
	if len(paths) > 0 {
		return paths
	}
	if m.Storage != nil && len(m.Storage.Items) > 0 && m.StorageCursor < len(m.Storage.Items) {
		return []string{m.Storage.Items[m.StorageCursor].Path}
	}
	return nil
}

func (m Model) handlePullSelected() (Model, tea.Cmd, bool) {
	paths := m.getMarkedOrCursorPaths()
	if len(paths) == 0 {
		m.Status = "No items selected to pull"
		return m, nil, true
	}
	destDir := storage.GetDefaultDownloadDir()
	m.Status = fmt.Sprintf("Pulling %d item(s) to %s...", len(paths), filepath.Base(destDir))
	return m, PullStorageItemCmd(m.Executor, paths, destDir), true
}

func (m Model) handleCopy() (Model, tea.Cmd, bool) {
	paths := m.getMarkedOrCursorPaths()
	if len(paths) == 0 {
		return m, nil, true
	}
	m.ClipboardPaths = paths
	m.ClipboardIsCut = false
	m.Status = fmt.Sprintf("Copied %d item(s) to clipboard. Press [p] in target folder to paste.", len(paths))
	return m, nil, true
}

func (m Model) handleCut() (Model, tea.Cmd, bool) {
	paths := m.getMarkedOrCursorPaths()
	if len(paths) == 0 {
		return m, nil, true
	}
	m.ClipboardPaths = paths
	m.ClipboardIsCut = true
	m.Status = fmt.Sprintf("Cut %d item(s) to clipboard. Press [p] in target folder to move.", len(paths))
	return m, nil, true
}

func (m Model) handlePaste() (Model, tea.Cmd, bool) {
	if len(m.ClipboardPaths) == 0 {
		m.Status = "Clipboard is empty. Copy [c] or cut [m] items first."
		return m, nil, true
	}
	targetDir := m.StoragePath
	if targetDir == "" {
		targetDir = "/sdcard"
	}
	items := m.ClipboardPaths
	isCut := m.ClipboardIsCut
	m.ClipboardPaths = nil
	m.SelectedFiles = make(map[string]bool)

	if isCut {
		m.Status = fmt.Sprintf("Moving %d item(s) to %s...", len(items), targetDir)
		return m, tea.Sequence(
			MoveStorageItemCmd(m.Executor, items, targetDir),
			LoadStorageCmd(m.Executor, targetDir),
		), true
	}

	m.Status = fmt.Sprintf("Copying %d item(s) to %s...", len(items), targetDir)
	return m, tea.Sequence(
		CopyStorageItemCmd(m.Executor, items, targetDir),
		LoadStorageCmd(m.Executor, targetDir),
	), true
}
