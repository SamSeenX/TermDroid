package storage

import (
	"os"
	"testing"
)

func TestGetDefaultDownloadDir(t *testing.T) {
	dir := GetDefaultDownloadDir()
	if dir == "" {
		t.Errorf("expected non-empty download directory")
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Errorf("expected download directory to exist and be a folder: %s", dir)
	}
}

func TestStorageOpsValidation(t *testing.T) {
	err := MakeDeviceDir(nil, "/sdcard", "")
	if err == nil {
		t.Errorf("expected error for empty folder name")
	}

	err = RenameDeviceItem(nil, "/sdcard/file.txt", "")
	if err == nil {
		t.Errorf("expected error for empty new name")
	}

	err = PushDeviceFile(nil, "", "/sdcard")
	if err == nil {
		t.Errorf("expected error for empty local path")
	}

	err = PushDeviceFile(nil, "/non/existent/path/termdroid_test.tmp", "/sdcard")
	if err == nil {
		t.Errorf("expected error for non-existent local file")
	}

	_, err = BatchPullDeviceFiles(nil, nil, "")
	if err == nil {
		t.Errorf("expected error for empty remote paths")
	}
}
