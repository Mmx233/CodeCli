package cmd

import (
	"os/exec"
	"runtime"
)

type ExecInterface interface {
	Command(name string, args ...string) error
	Background(name string, args ...string) error
	SetDir(dir string) ExecInterface
}

var Exec ExecInterface

type defaultExec struct {
	dir string
}

func (a defaultExec) Command(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = a.dir
	return cmd.Run()
}
func (a defaultExec) Background(name string, args ...string) error {
	return a.Command(name, args...)
}
func (a defaultExec) SetDir(dir string) ExecInterface {
	return defaultExec{dir: dir}
}

type windowsExec struct {
	dir string
}

func (a windowsExec) Command(name string, args ...string) error {
	cmd := exec.Command("cmd", append([]string{
		"/c", "start", name,
	}, args...)...)
	cmd.Dir = a.dir
	return cmd.Start()
}
func (a windowsExec) Background(name string, args ...string) error {
	return a.Command("/B", append([]string{name}, args...)...)
}
func (a windowsExec) SetDir(dir string) ExecInterface {
	return windowsExec{dir: dir}
}

func init() {
	switch runtime.GOOS {
	case "windows":
		Exec = windowsExec{}
	default:
		Exec = defaultExec{}
	}
}

func OpenIdea(idea, projectPath string) error {
	return Exec.Background(idea, projectPath)
}

func RunCmd(path, program string) error {
	return Exec.SetDir(path).Command(program)
}
