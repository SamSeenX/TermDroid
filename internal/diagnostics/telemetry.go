package diagnostics

import (
	"fmt"
	"strconv"
	"strings"
	"termdroid/internal/bridge"
	"time"
)

type DeviceVitals struct {
	BatteryLevel       string
	BatteryTemperature string
	BatteryHealth      string
	BatteryCharging    bool
	CPULoad            string
	RAMSummary         string
}

func FetchVitals(exec *bridge.Executor) (*DeviceVitals, error) {
	vitals := &DeviceVitals{
		BatteryLevel:       "N/A",
		BatteryTemperature: "N/A",
		BatteryHealth:      "Good",
		BatteryCharging:    false,
		CPULoad:            "Normal",
		RAMSummary:         "N/A",
	}

	batOut, err := exec.RunADBWithTimeout(3*time.Second, "shell", "dumpsys", "battery")
	if err == nil {
		for _, line := range strings.Split(batOut, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "level:") {
				vitals.BatteryLevel = strings.TrimSpace(strings.TrimPrefix(line, "level:")) + "%"
			} else if strings.HasPrefix(line, "temperature:") {
				tempStr := strings.TrimSpace(strings.TrimPrefix(line, "temperature:"))
				if t, parseErr := strconv.Atoi(tempStr); parseErr == nil {
					vitals.BatteryTemperature = fmt.Sprintf("%.1f C", float64(t)/10.0)
				}
			} else if strings.HasPrefix(line, "AC powered: true") || strings.HasPrefix(line, "USB powered: true") {
				vitals.BatteryCharging = true
			}
		}
	}

	cpuOut, err := exec.RunADBWithTimeout(3*time.Second, "shell", "dumpsys", "cpuinfo")
	if err == nil {
		lines := strings.Split(cpuOut, "\n")
		for _, line := range lines {
			if strings.Contains(line, "TOTAL:") {
				vitals.CPULoad = strings.TrimSpace(line)
				break
			}
		}
	}

	return vitals, nil
}
