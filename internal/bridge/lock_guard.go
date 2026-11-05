package bridge

import (
	"strings"
	"time"
)

type ScreenLockState struct {
	IsAwake  bool
	IsLocked bool
}

func (e *Executor) CheckLockState() ScreenLockState {
	state := ScreenLockState{
		IsAwake:  true,
		IsLocked: false,
	}

	powerOut, err := e.RunADBWithTimeout(2*time.Second, "shell", "dumpsys", "power")
	if err == nil {
		if strings.Contains(powerOut, "mWakefulness=Asleep") ||
			strings.Contains(powerOut, "mWakefulness=Dozing") {
			state.IsAwake = false
		}
	}

	windowOut, err := e.RunADBWithTimeout(2*time.Second, "shell", "dumpsys", "window")
	if err == nil {
		if strings.Contains(windowOut, "mDreamingLockscreen=true") ||
			strings.Contains(windowOut, "isStatusBarKeyguard=true") ||
			strings.Contains(windowOut, "mShowingDream=true") ||
			strings.Contains(windowOut, "mKeyguardOccluded=false") && strings.Contains(windowOut, "Keyguard") {
			state.IsLocked = true
		}
	}

	if !state.IsAwake {
		state.IsLocked = true
	}

	return state
}
