package idea

import (
	"golang.org/x/sys/windows"
	"os/exec"
	"syscall"
)

func init() {
	Exec = windowsExec{}
	ideaCommand = func(idea string) string {
		return idea + "64.exe"
	}
}

type windowsExec struct {
	dir string
}

func (a windowsExec) Command(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = a.dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:       false,
		CreationFlags:    windows.CREATE_NEW_CONSOLE,
		NoInheritHandles: true,
	}
	return cmd.Start()
}

func (a windowsExec) CreateProcess(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = a.dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:       true,
		CreationFlags:    windows.CREATE_NEW_CONSOLE | windows.CREATE_NO_WINDOW,
		NoInheritHandles: true,
	}
	return cmd.Start()
}
func (a windowsExec) SetDir(dir string) ExecInterface {
	return windowsExec{dir: dir}
}
