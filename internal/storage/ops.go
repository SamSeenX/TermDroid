package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"termdroid/internal/bridge"
	"time"
)

func GetDefaultDownloadDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	downloads := filepath.Join(home, "Downloads")
	if info, err := os.Stat(downloads); err == nil && info.IsDir() {
		return downloads
	}
	return home
}

func PushDeviceFile(exec *bridge.Executor, localPath, remoteDir string) error {
	localPath = strings.TrimSpace(localPath)
	if localPath == "" {
		return errors.New("local path cannot be empty")
	}
	if strings.HasPrefix(localPath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			localPath = filepath.Join(home, localPath[2:])
		}
	}
	if _, err := os.Stat(localPath); err != nil {
		return fmt.Errorf("local file does not exist: %w", err)
	}
	_, err := exec.RunADBWithTimeout(120*time.Second, "push", localPath, remoteDir)
	return err
}

func CopyDeviceItem(exec *bridge.Executor, srcPath, destDir string) error {
	srcQuoted := fmt.Sprintf("\"%s\"", strings.TrimSpace(srcPath))
	destQuoted := fmt.Sprintf("\"%s\"", strings.TrimSpace(destDir))
	_, err := exec.RunADBWithTimeout(30*time.Second, "shell", "cp", "-r", srcQuoted, destQuoted)
	return err
}

func MoveDeviceItem(exec *bridge.Executor, srcPath, destDir string) error {
	srcQuoted := fmt.Sprintf("\"%s\"", strings.TrimSpace(srcPath))
	destQuoted := fmt.Sprintf("\"%s\"", strings.TrimSpace(destDir))
	_, err := exec.RunADBWithTimeout(15*time.Second, "shell", "mv", srcQuoted, destQuoted)
	return err
}

func MakeDeviceDir(exec *bridge.Executor, parentDir, folderName string) error {
	folderName = strings.TrimSpace(folderName)
	if folderName == "" {
		return errors.New("folder name cannot be empty")
	}
	target := filepath.ToSlash(filepath.Join(parentDir, folderName))
	targetQuoted := fmt.Sprintf("\"%s\"", target)
	_, err := exec.RunADBWithTimeout(10*time.Second, "shell", "mkdir", "-p", targetQuoted)
	return err
}

func RenameDeviceItem(exec *bridge.Executor, oldPath, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("new name cannot be empty")
	}
	dir := filepath.Dir(oldPath)
	newPath := filepath.ToSlash(filepath.Join(dir, newName))
	return MoveDeviceItem(exec, oldPath, newPath)
}

func BatchPullDeviceFiles(exec *bridge.Executor, remotePaths []string, localDir string) (int, error) {
	if len(remotePaths) == 0 {
		return 0, errors.New("no files selected to pull")
	}
	if localDir == "" {
		localDir = GetDefaultDownloadDir()
	}
	successCount := 0
	for _, remote := range remotePaths {
		if strings.TrimSpace(remote) == "" {
			continue
		}
		_, err := exec.RunADBWithTimeout(60*time.Second, "pull", remote, localDir)
		if err == nil {
			successCount++
		}
	}
	if successCount == 0 && len(remotePaths) > 0 {
		return 0, errors.New("failed to pull selected files")
	}
	return successCount, nil
}
