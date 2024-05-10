package idea

import (
	"golang.org/x/sys/windows"
	"os/exec"
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

func (a windowsExec) Command(name string, args ...string) error {
	var argList = []string{
		"-NoProfile", "-Command", "Start-Process",
		name,
	}
	if len(args) != 0 {
		argList = append(argList, "-ArgumentList")
		argList = append(argList, args...)
	}
	cmd := exec.Command("powershell", argList...)
	cmd.Dir = a.dir
	return cmd.Run()
}

func (a windowsExec) CreateProcess(name string, args ...string) error {
	program := name
	if len(args) != 0 {
		program += " " + strings.Join(args, " ")
	}
	program = strings.Replace(program, `/`, `\\`, -1)

	var workdir *uint16
	if a.dir != "" {
		workdir = windows.StringToUTF16Ptr(strings.Replace(a.dir, `/`, `\\`, -1))
	}

	var procInfo syscall.ProcessInformation
	startupInfo := &syscall.StartupInfo{
		StdErr:    syscall.Stderr,
		StdOutput: syscall.Stdout,
		StdInput:  syscall.Stdin,
	}
	return syscall.CreateProcess(
		nil, windows.StringToUTF16Ptr(program),
		nil, nil, false, 0, nil,
		workdir, startupInfo, &procInfo)
}
func (a windowsExec) SetDir(dir string) ExecInterface {
	return windowsExec{dir: dir}
}
