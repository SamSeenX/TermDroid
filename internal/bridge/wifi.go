package bridge

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (e *Executor) GetDeviceIP() (string, error) {
	out, err := e.RunADBWithTimeout(3*time.Second, "shell", "ip", "route")
	if err == nil && out != "" {
		for _, line := range strings.Split(out, "\n") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "src" && i+1 < len(fields) {
					return fields[i+1], nil
				}
			}
		}
	}

	outWlan, errWlan := e.RunADBWithTimeout(3*time.Second, "shell", "ip", "-f", "inet", "addr", "show", "wlan0")
	if errWlan == nil && outWlan != "" {
		for _, line := range strings.Split(outWlan, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "inet ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					cidr := parts[1]
					ip := strings.Split(cidr, "/")[0]
					return ip, nil
				}
			}
		}
	}

	return "", errors.New("unable to resolve device local wifi ip")
}

func (e *Executor) SwitchToWireless(port string) (string, error) {
	if port == "" {
		port = "5555"
	}

	ip, err := e.GetDeviceIP()
	if err != nil {
		return "", err
	}

	_, err = e.RunADBWithTimeout(5*time.Second, "tcpip", port)
	if err != nil {
		return "", fmt.Errorf("failed to enable tcpip mode: %w", err)
	}

	time.Sleep(1500 * time.Millisecond)

	target := fmt.Sprintf("%s:%s", ip, port)
	out, err := e.RunADBWithTimeout(6*time.Second, "connect", target)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", target, err)
	}

	if strings.Contains(out, "connected") {
		return target, nil
	}

	return target, fmt.Errorf("connection response unexpected: %s", out)
}
