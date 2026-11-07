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

func CleanHostPath(raw string) string {
	raw = strings.TrimSpace(raw)
	if (strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'")) ||
		(strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"")) {
		if len(raw) >= 2 {
			raw = raw[1 : len(raw)-1]
		}
	}
	raw = strings.ReplaceAll(raw, "\\ ", " ")
	raw = strings.ReplaceAll(raw, "\\(", "(")
	raw = strings.ReplaceAll(raw, "\\)", ")")
	if strings.HasPrefix(raw, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			raw = filepath.Join(home, raw[2:])
		}
	}
	return strings.TrimSpace(raw)
}

func ResolveUploadDestination(remoteDir string) string {
	remoteDir = strings.TrimSpace(remoteDir)
	if remoteDir == "" || remoteDir == "/sdcard" || remoteDir == "/sdcard/" {
		return "/sdcard/Download/TermDroid"
	}
	return remoteDir
}

func PushDeviceFile(exec *bridge.Executor, localPath, remoteDir string) (string, error) {
	localPath = CleanHostPath(localPath)
	if localPath == "" {
		return "", errors.New("local path cannot be empty")
	}
	if _, err := os.Stat(localPath); err != nil {
		return "", fmt.Errorf("local file does not exist: %w", err)
	}

	targetDir := ResolveUploadDestination(remoteDir)
	mkdirQuoted := fmt.Sprintf("\"%s\"", targetDir)
	exec.RunADBWithTimeout(5*time.Second, "shell", "mkdir", "-p", mkdirQuoted)

	_, err := exec.RunADBWithTimeout(120*time.Second, "push", localPath, targetDir)
	if err != nil {
		return "", err
	}

	destFile := filepath.ToSlash(filepath.Join(targetDir, filepath.Base(localPath)))
	TriggerMediaScan(exec, destFile)

	return targetDir, nil
}

func TriggerMediaScan(exec *bridge.Executor, filePath string) {
	fileUri := fmt.Sprintf("'file://%s'", strings.Trim(filePath, "'\""))
	exec.RunADBWithTimeout(3*time.Second, "shell", "am", "broadcast", "-a", "android.intent.action.MEDIA_SCANNER_SCAN_FILE", "-d", fileUri)
	exec.RunADBWithTimeout(3*time.Second, "shell", "content", "call", "--uri", "content://media", "--method", "scan_volume", "--extra", "volume:s:external_primary")
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
