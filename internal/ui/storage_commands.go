package ui

import (
	"fmt"
	"termdroid/internal/bridge"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func DeleteStorageItemCmd(exec *bridge.Executor, targetPath string) tea.Cmd {
	return func() tea.Msg {
		err := storage.DeleteDevicePath(exec, targetPath)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Delete failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Deleted: %s", targetPath))
	}
}
