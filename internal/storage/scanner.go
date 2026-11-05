package storage

import (
	"sort"
	"strconv"
	"strings"
	"termdroid/internal/bridge"
	"time"
)

type StorageSummary struct {
	TotalSpace string
	UsedSpace  string
	FreeSpace  string
	Percent    string
	Items      []StorageItem
}

func ScanStorageTree(exec *bridge.Executor, rootPath string) (*StorageSummary, error) {
	if rootPath == "" {
		rootPath = "/sdcard"
	}

	summary := &StorageSummary{
		TotalSpace: "Unknown",
		UsedSpace:  "Unknown",
		FreeSpace:  "Unknown",
		Percent:    "0%",
	}

	dfOut, err := exec.RunADBWithTimeout(4*time.Second, "shell", "df", "-h", rootPath)
	if err == nil {
		lines := strings.Split(dfOut, "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				summary.TotalSpace = fields[1]
				summary.UsedSpace = fields[2]
				summary.FreeSpace = fields[3]
				summary.Percent = fields[4]
			}
		}
	}

	duScript := "for d in /sdcard/*; do [ -d \"$d\" ] && [ \"$d\" != \"/sdcard/Android\" ] && du -sk \"$d\" 2>/dev/null; done; true"
	duOut, err := exec.RunADBWithTimeout(8*time.Second, "shell", duScript)
	if err != nil || duOut == "" {
		duOut, _ = exec.RunADBWithTimeout(8*time.Second, "shell", "du -sk /sdcard/DCIM /sdcard/Download /sdcard/Movies /sdcard/Pictures /sdcard/Music /sdcard/Documents 2>/dev/null; true")
	}

	var items []StorageItem
	for _, line := range strings.Split(duOut, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		kb, parseErr := strconv.ParseInt(fields[0], 10, 64)
		if parseErr != nil {
			continue
		}

		path := fields[1]
		if path == rootPath || strings.TrimSuffix(path, "/") == strings.TrimSuffix(rootPath, "/") {
			continue
		}

		bytes := kb * 1024
		items = append(items, StorageItem{
			Path:      path,
			SizeBytes: bytes,
			SizeHuman: FormatBytes(bytes),
			IsDir:     true,
			Category:  CategorizePath(path),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SizeBytes > items[j].SizeBytes
	})

	summary.Items = items
	return summary, nil
}
