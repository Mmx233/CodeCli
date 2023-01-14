package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/CodeCli/global"
)

func Open(addr string) error {
	project, e := LoadProject(addr)
	if e != nil {
		return e
	}
	return OpenProject(project.Path)
}

func OpenProject(path string) error {
	var idea = global.Commands.Project.Idea
	if idea == "" {
		var e error
		idea, e = IdeaSelect(path)
		if e != nil {
			return e
		}
	}
	return cmd.OpenIdea(idea, path)
}

func OpenCmd(addr string) error {
	project, e := LoadProject(addr)
	if e != nil {
		return e
	}
	return cmd.RunCmd(project.Path, global.Config.Default.CmdProgram)
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
