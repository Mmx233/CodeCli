package desktop

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func setProcessAttributes(c *exec.Cmd, console bool) {
	flags := uint32(windows.CREATE_NEW_PROCESS_GROUP)
	if console {
		flags |= windows.CREATE_NEW_CONSOLE
	} else {
		flags |= windows.CREATE_NO_WINDOW
	}
	c.SysProcAttr = &syscall.SysProcAttr{CreationFlags: flags, HideWindow: !console, NoInheritHandles: true}
}
