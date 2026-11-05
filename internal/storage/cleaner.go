package storage

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"termdroid/internal/bridge"
	"time"
)

func FindLargeFiles(exec *bridge.Executor, minSizeMB int) ([]StorageItem, error) {
	if minSizeMB <= 0 {
		minSizeMB = 50
	}

	arg := fmt.Sprintf("+%dM", minSizeMB)
	cmdOut, err := exec.RunADBWithTimeout(15*time.Second, "shell", "find", "/sdcard", "-type", "f", "-size", arg, "-exec", "ls", "-s", "{}", "+")
	if err != nil {
		cmdOut, err = exec.RunADBWithTimeout(15*time.Second, "shell", "find", "/sdcard/Download", "-type", "f", "-size", arg)
	}

	var items []StorageItem
	for _, line := range strings.Split(cmdOut, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		var path string
		var bytes int64

		if len(fields) >= 2 {
			blocks, err := strconv.ParseInt(fields[0], 10, 64)
			if err == nil {
				bytes = blocks * 1024
				path = fields[1]
			} else {
				path = fields[0]
			}
		} else {
			path = fields[0]
		}

		if path == "" {
			continue
		}

		if bytes == 0 {
			bytes = int64(minSizeMB) * 1024 * 1024
		}

		items = append(items, StorageItem{
			Path:      path,
			SizeBytes: bytes,
			SizeHuman: FormatBytes(bytes),
			IsDir:     false,
			Category:  CategorizePath(path),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SizeBytes > items[j].SizeBytes
	})

	return items, nil
}

func DeleteDevicePath(exec *bridge.Executor, targetPath string) error {
	_, err := exec.RunADBWithTimeout(8*time.Second, "shell", "rm", "-rf", targetPath)
	return err
}

func PullDeviceFile(exec *bridge.Executor, remotePath, localPath string) error {
	_, err := exec.RunADBWithTimeout(30*time.Second, "pull", remotePath, localPath)
	return err
}
