package cmd

import (
	"os/exec"
	"runtime"
)

type (
	ExecInterface interface {
		Command(name string, args ...string) error
		Background(name string, args ...string) error
		SetDir(dir string) ExecInterface
	}
)

var (
	ideaCommand = func(idea string) string {
		return idea
	}
	Exec ExecInterface
)

type defaultExec struct {
	dir string
}

func (a defaultExec) Command(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = a.dir
	return cmd.Run()
}
func (a defaultExec) Background(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = a.dir
	return cmd.Run()
}
func (a defaultExec) SetDir(dir string) ExecInterface {
	return defaultExec{dir: dir}
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
func (a windowsExec) Background(name string, args ...string) error {
	var argList = []string{
		"-WindowStyle", "Hidden", "-NoProfile", "-Command", "Start-Process",
		name,
	}
	if len(args) != 0 {
		argList = append(argList, "-ArgumentList")
		argList = append(argList, args...)
	}
	argList = append(argList, "-WindowStyle", "Hidden")
	cmd := exec.Command("powershell", argList...)
	cmd.Dir = a.dir
	return cmd.Start()
}
func (a windowsExec) SetDir(dir string) ExecInterface {
	return windowsExec{dir: dir}
}

func init() {
	switch runtime.GOOS {
	case "windows":
		Exec = windowsExec{}
		ideaCommand = func(idea string) string {
			return idea + "64.exe"
		}
	default:
		Exec = defaultExec{}
	}
}

func OpenIdea(idea, projectPath string) error {
	return Exec.Background(ideaCommand(idea), projectPath)
}

func RunCmd(path, program string) error {
	return Exec.SetDir(path).Command(program)
}
