package project

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/CodeCli/global"
)

func Open(addr string) error {
	_, _, projectPath, e := PrepareProjectFiles(addr)
	if e != nil {
		return e
	}
	return OpenProject(projectPath)
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
	_, _, projectPath, e := PrepareProjectFiles(addr)
	if e != nil {
		return e
	}
	return cmd.RunCmd(projectPath, global.Config.Default.CmdProgram)
}
