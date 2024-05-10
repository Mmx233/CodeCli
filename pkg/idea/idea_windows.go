package idea

import (
	"golang.org/x/sys/windows"
	"strings"
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

func (a windowsExec) program(name string, args ...string) (cmdLine *uint16, workdir *uint16, err error) {
	program := name
	if len(args) != 0 {
		program += " " + strings.Join(args, " ")
	}
	program = strings.Replace(program, `/`, `\\`, -1)
	cmdLine, err = windows.UTF16PtrFromString(program)
	if err != nil {
		return
	}

	if a.dir != "" {
		workdir = windows.StringToUTF16Ptr(strings.Replace(a.dir, `/`, `\\`, -1))
	}
	return
}

func (a windowsExec) Command(name string, args ...string) error {
	program, workdir, err := a.program(name, args...)
	if err != nil {
		return err
	}

	var procInfo syscall.ProcessInformation
	startupInfo := &syscall.StartupInfo{
		Flags:      windows.STARTF_USESHOWWINDOW,
		ShowWindow: windows.SW_NORMAL,
	}
	return syscall.CreateProcess(
		nil, program,
		nil, nil, false, windows.CREATE_NEW_CONSOLE, nil,
		workdir, startupInfo, &procInfo)
}

func (a windowsExec) CreateProcess(name string, args ...string) error {
	program, workdir, err := a.program(name, args...)
	if err != nil {
		return err
	}

	var procInfo syscall.ProcessInformation
	startupInfo := &syscall.StartupInfo{}
	return syscall.CreateProcess(
		nil, program,
		nil, nil, false, windows.CREATE_NO_WINDOW, nil,
		workdir, startupInfo, &procInfo)
}
func (a windowsExec) SetDir(dir string) ExecInterface {
	return windowsExec{dir: dir}
}
