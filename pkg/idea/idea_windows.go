package idea

import "os/exec"

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
