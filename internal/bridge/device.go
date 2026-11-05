package bridge

import (
	"strings"
	"time"
)

type DeviceInfo struct {
	ID            string
	State         string
	Model         string
	Product       string
	IsWireless    bool
	IPAddress     string
	BatteryLevel  string
	BatteryStatus string
}

func (e *Executor) ListDevices() ([]DeviceInfo, error) {
	out, err := e.RunADBWithTimeout(4*time.Second, "devices", "-l")
	if err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		devID := fields[0]
		state := fields[1]
		if state != "device" {
			continue
		}

		dev := DeviceInfo{
			ID:         devID,
			State:      state,
			IsWireless: strings.Contains(devID, ":"),
		}

		if dev.IsWireless {
			parts := strings.Split(devID, ":")
			dev.IPAddress = parts[0]
		}

		for _, field := range fields[2:] {
			if strings.HasPrefix(field, "model:") {
				dev.Model = strings.TrimPrefix(field, "model:")
			} else if strings.HasPrefix(field, "product:") {
				dev.Product = strings.TrimPrefix(field, "product:")
			}
		}

		if dev.Model == "" {
			dev.Model = "Android Device"
		}

		devices = append(devices, dev)
	}

	return devices, nil
}
