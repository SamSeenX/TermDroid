package diagnostics

import (
	"strings"
	"termdroid/internal/bridge"
	"time"
)

type ProcessItem struct {
	User string
	PID  string
	PPID string
	VSZ  string
	RSS  string
	Name string
}

func ListProcesses(exec *bridge.Executor) ([]ProcessItem, error) {
	out, err := exec.RunADBWithTimeout(5*time.Second, "shell", "ps", "-A")
	if err != nil {
		out, err = exec.RunADBWithTimeout(5*time.Second, "shell", "ps")
	}
	if err != nil {
		return nil, err
	}

	var procs []ProcessItem
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		user := fields[0]
		pid := fields[1]
		ppid := fields[2]
		vsz := fields[3]
		rss := fields[4]
		name := fields[len(fields)-1]

		if name == "ps" || name == "sh" {
			continue
		}

		procs = append(procs, ProcessItem{
			User: user,
			PID:  pid,
			PPID: ppid,
			VSZ:  vsz,
			RSS:  rss,
			Name: name,
		})
	}

	return procs, nil
}

func KillProcess(exec *bridge.Executor, pid string) error {
	_, err := exec.RunADBWithTimeout(3*time.Second, "shell", "kill", "-9", pid)
	return err
}

func ForceStopPackage(exec *bridge.Executor, pkgName string) error {
	_, err := exec.RunADBWithTimeout(3*time.Second, "shell", "am", "force-stop", pkgName)
	return err
}
