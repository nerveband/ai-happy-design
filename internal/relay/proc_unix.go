//go:build !windows

package relay

import (
	"os/exec"
	"syscall"
)

// detachProcess ensures background relay keeps running after parent exits.
func detachProcess(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
