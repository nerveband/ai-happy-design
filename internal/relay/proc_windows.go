//go:build windows

package relay

import "os/exec"

// detachProcess is a no-op on Windows.
func detachProcess(_ *exec.Cmd) {}
