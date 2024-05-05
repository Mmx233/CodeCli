package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/idea"
)

func Open(addr string) error {
	project, err := LoadProject(addr)
	if err != nil {
		return err
	}
	return OpenProject(project.Path)
}

func OpenProject(path string) error {
	var ideaName = global.Commands.Project.Idea
	if ideaName == "" {
		var err error
		ideaName, err = IdeaSelect(path)
		if err != nil {
			return err
		}
	}
	return idea.Open(ideaName, path)
}

func OpenCmd(addr string) error {
	project, err := LoadProject(addr)
	if err != nil {
		return err
	}
	return idea.RunCmd(project.Path, global.Config.Default.CmdProgram)
}

type Project struct {
	GitSite  string
	Username string
	Repo     string
	SubDir   string

	// dir of repos
	Dir  string
	Path string
}

func (a Project) Url() string {
	return fmt.Sprintf("https://%s/%s/%s.git", a.GitSite, a.Username, a.Repo)
}

func (a Project) Open() error {
	return Open(a.Url())
}

func (a Project) OpenCmd() error {
	return OpenCmd(a.Url())
}
