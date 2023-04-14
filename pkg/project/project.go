package project

import (
	"fmt"
	global2 "github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/idea"
)

func Open(addr string) error {
	project, e := LoadProject(addr)
	if e != nil {
		return e
	}
	return OpenProject(project.Path)
}

func OpenProject(path string) error {
	var ideaName = global2.Commands.Project.Idea
	if ideaName == "" {
		var e error
		ideaName, e = IdeaSelect(path)
		if e != nil {
			return e
		}
	}
	return idea.Open(ideaName, path)
}

func OpenCmd(addr string) error {
	project, e := LoadProject(addr)
	if e != nil {
		return e
	}
	return idea.RunCmd(project.Path, global2.Config.Default.CmdProgram)
}

type Project struct {
	GitSite  string
	Username string
	Repo     string

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
