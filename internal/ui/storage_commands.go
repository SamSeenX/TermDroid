package ui

import (
	"fmt"
	"path/filepath"
	"termdroid/internal/bridge"
	"termdroid/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type StorageModalAction int

const (
	ModalNone StorageModalAction = iota
	ModalMkdir
	ModalRename
	ModalPush
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

func BatchDeleteStorageCmd(exec *bridge.Executor, targets []string) tea.Cmd {
	return func() tea.Msg {
		deleted := 0
		for _, target := range targets {
			if storage.DeleteDevicePath(exec, target) == nil {
				deleted++
			}
		}
		return StatusMsg(fmt.Sprintf("Deleted %d items from device", deleted))
	}
}

func PullStorageItemCmd(exec *bridge.Executor, remotePaths []string, localDir string) tea.Cmd {
	return func() tea.Msg {
		if localDir == "" {
			localDir = storage.GetDefaultDownloadDir()
		}
		count, err := storage.BatchPullDeviceFiles(exec, remotePaths, localDir)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Pull failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Pulled %d file(s) to %s", count, filepath.Base(localDir)))
	}
}

func PushStorageItemCmd(exec *bridge.Executor, localPath, remoteDir string) tea.Cmd {
	return func() tea.Msg {
		targetDir, err := storage.PushDeviceFile(exec, localPath, remoteDir)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Push failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Uploaded %s to %s", filepath.Base(localPath), targetDir))
	}
}

func CopyStorageItemCmd(exec *bridge.Executor, srcPaths []string, destDir string) tea.Cmd {
	return func() tea.Msg {
		success := 0
		for _, src := range srcPaths {
			if storage.CopyDeviceItem(exec, src, destDir) == nil {
				success++
			}
		}
		return StatusMsg(fmt.Sprintf("Copied %d item(s) to %s", success, filepath.Base(destDir)))
	}
}

func MoveStorageItemCmd(exec *bridge.Executor, srcPaths []string, destDir string) tea.Cmd {
	return func() tea.Msg {
		success := 0
		for _, src := range srcPaths {
			if storage.MoveDeviceItem(exec, src, destDir) == nil {
				success++
			}
		}
		return StatusMsg(fmt.Sprintf("Moved %d item(s) to %s", success, filepath.Base(destDir)))
	}
}

func MkdirStorageItemCmd(exec *bridge.Executor, parentDir, folderName string) tea.Cmd {
	return func() tea.Msg {
		err := storage.MakeDeviceDir(exec, parentDir, folderName)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Mkdir failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Created folder: %s", folderName))
	}
}

func RenameStorageItemCmd(exec *bridge.Executor, oldPath, newName string) tea.Cmd {
	return func() tea.Msg {
		err := storage.RenameDeviceItem(exec, oldPath, newName)
		if err != nil {
			return StatusMsg(fmt.Sprintf("Rename failed: %v", err))
		}
		return StatusMsg(fmt.Sprintf("Renamed to: %s", newName))
	}
}
