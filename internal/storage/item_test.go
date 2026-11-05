package storage

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, c := range cases {
		actual := FormatBytes(c.input)
		if actual != c.expected {
			t.Errorf("FormatBytes(%d) = %s, expected %s", c.input, actual, c.expected)
		}
	}
}

func TestCategorizePath(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{"/sdcard/DCIM/Camera", "Photos & Media"},
		{"/sdcard/Android/media/com.whatsapp", "Messaging Media"},
		{"/sdcard/Download/installer.apk", "Downloads"},
		{"/sdcard/Music/track.mp3", "Audio & Music"},
		{"/sdcard/random_folder", "Other Storage"},
	}

	for _, c := range cases {
		actual := CategorizePath(c.path)
		if actual != c.expected {
			t.Errorf("CategorizePath(%s) = %s, expected %s", c.path, actual, c.expected)
		}
	}
}
