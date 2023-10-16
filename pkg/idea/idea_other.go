//go:build !windows

package idea

import "os/exec"

func init() {
	Exec = defaultExec{}
}

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
