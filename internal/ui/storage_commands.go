package ui

import (
	"fmt"
	"termdroid/internal/bridge"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func DeleteStorageItemCmd(exec *bridge.Executor, targetPath, parentPath string) tea.Cmd {
	return func() tea.Msg {
		err := storage.DeleteDevicePath(exec, targetPath)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Delete failed: %v", err))
		}
		newSummary, scanErr := storage.ScanStorageTree(exec, parentPath)
		if scanErr != nil {
			return StatusMsg(fmt.Sprintf("Deleted %s (error refreshing list)", targetPath))
		}
		return StorageLoadedMsg(newSummary)
	}
}
