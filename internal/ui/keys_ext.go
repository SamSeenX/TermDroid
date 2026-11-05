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

	if m.StorageCursor < len(m.Storage.Items) {
		selected := m.Storage.Items[m.StorageCursor]
		m.DeleteTarget = selected.Path
		m.ConfirmDelete = true
		m.Status = fmt.Sprintf("Delete %s? Press [y] to confirm, [Esc] to cancel", selected.Path)
	}
	return m, nil
}

func (m Model) handleConfirmDelete() (tea.Model, tea.Cmd) {
	if m.ConfirmDelete && m.DeleteTarget != "" {
		target := m.DeleteTarget
		m.ConfirmDelete = false
		m.DeleteTarget = ""

		var remaining []storage.StorageItem
		for _, it := range m.Storage.Items {
			if it.Path != target {
				remaining = append(remaining, it)
			}
		}
		m.Storage.Items = remaining

		if m.StorageCursor >= len(remaining) && len(remaining) > 0 {
			m.StorageCursor = len(remaining) - 1
		}

		m.Status = fmt.Sprintf("Deleting %s from device...", target)
		return m, DeleteStorageItemCmd(m.Executor, target)
	}
	return m, nil
}

func (m Model) handleActionShortcut(keyStr string) (tea.Model, tea.Cmd) {
	if act, ok := findActionByShortcut(keyStr); ok {
		m.Status = fmt.Sprintf("Shortcut fired: %s...", act.Name)
		return m, ExecuteKeypadCmd(m.Executor, act)
	}
	return m, nil
}
