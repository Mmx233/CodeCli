//go:build !windows

package idea

import (
	"os/exec"
)

func init() {
	Exec = defaultExec{}
}

type defaultExec struct {
	dir string
}

func (a defaultExec) NewWindowCommand(name string, args ...string) error {
	return a.CreateProcess(name, args...)
}

func (a defaultExec) CreateProcess(name string, args ...string) error {
	cmd := exec.Command("nohup", append([]string{name}, args...)...)
	cmd.Dir = a.dir
	return cmd.Start()
}

func (a defaultExec) SetDir(dir string) ExecInterface {
	return defaultExec{dir: dir}
}
