package cmd

import (
	"os/exec"
	"runtime"
	"syscall"
)

var ideaExec func(name string, args ...string) error

func defaultExec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd.Start()
}
func windowsExec(name string, args ...string) error {
	return exec.Command("cmd", append([]string{
		"/c", "start", "/B",
		name + "64.exe",
	}, args...)...).Start()
}

func init() {
	switch runtime.GOOS {
	case "windows":
		ideaExec = windowsExec
	default:
		ideaExec = defaultExec
	}
}

func Goland(path string) error {
	return ideaExec("goland", path)
}
func Webstorm(path string) error {
	return ideaExec("webstorm", path)
}
