//go:build !windows

package desktop

import (
	"os/exec"
	"syscall"
)

func setProcessAttributes(c *exec.Cmd, console bool) {
	c.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
