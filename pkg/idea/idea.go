package idea

type (
	ExecInterface interface {
		NewWindowCommand(name string, args ...string) error
		CreateProcess(name string, args ...string) error
		SetDir(dir string) ExecInterface
	}
)

var (
	ideaCommand = func(idea string) string {
		return idea
	}
	Exec ExecInterface
)

func Open(idea, projectPath string) error {
	return Exec.CreateProcess(ideaCommand(idea), projectPath)
}

func RunCmd(path, program string) error {
	return Exec.SetDir(path).NewWindowCommand(program)
}
