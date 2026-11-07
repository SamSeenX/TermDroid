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

	_, err = PushDeviceFile(nil, "", "/sdcard")
	if err == nil {
		t.Errorf("expected error for empty local path")
	}

	_, err = PushDeviceFile(nil, "/non/existent/path/termdroid_test.tmp", "/sdcard")
	if err == nil {
		t.Errorf("expected error for non-existent local file")
	}

	_, err = BatchPullDeviceFiles(nil, nil, "")
	if err == nil {
		t.Errorf("expected error for empty remote paths")
	}
}

func TestResolveUploadDestination(t *testing.T) {
	if got := ResolveUploadDestination(""); got != "/sdcard/Download/TermDroid" {
		t.Errorf("expected /sdcard/Download/TermDroid, got %s", got)
	}
	if got := ResolveUploadDestination("/sdcard"); got != "/sdcard/Download/TermDroid" {
		t.Errorf("expected /sdcard/Download/TermDroid, got %s", got)
	}
	if got := ResolveUploadDestination("/sdcard/DCIM"); got != "/sdcard/DCIM" {
		t.Errorf("expected /sdcard/DCIM, got %s", got)
	}
}

func TestCleanHostPath(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"'path/to/file.jpg'", "path/to/file.jpg"},
		{"\"path/to/file.jpg\"", "path/to/file.jpg"},
		{"path/to/my\\ file.jpg", "path/to/my file.jpg"},
		{"path/to/flyer\\ \\(A4\\).png", "path/to/flyer (A4).png"},
	}

	for _, c := range cases {
		if got := CleanHostPath(c.input); got != c.expected {
			t.Errorf("CleanHostPath(%s) = %s, expected %s", c.input, got, c.expected)
		}
	}
}
