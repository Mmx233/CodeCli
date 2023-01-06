package project

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/CodeCli/util"
	"github.com/Mmx233/tool"
	"os"
	"strings"
)

func Open(addr string) error {
	switch len(strings.Split(addr, "/")) {
	case 1:
		if global.Config.Default.Username == "" {
			return util.ErrEmptyDefaultUsername
		}
		addr = global.Config.Default.Username + "/" + addr
		fallthrough
	case 2:
		if global.Config.Default.GitSite == "" {
			return util.ErrEmptyDefaultGitSite
		}
		addr = global.Config.Default.GitSite + "/" + addr
	case 3:
		break
	default:
		return util.ErrUnknownInput
	}

	split := strings.Split(addr, "/")
	projectDir := file.JoinPath(append([]string{global.Config.Storage.ProjectDir}, split[:2]...)...)

	e := os.MkdirAll(projectDir, 0600)
	if e != nil {
		return e
	}

	projectPath := file.JoinPath(projectDir, split[2])

	if !tool.File.Exists(projectPath) {
		if e = cmd.Clone("https://"+addr+".git", projectPath); e != nil {
			return e
		}
	}

	return OpenProject(projectPath)
}

func OpenProject(path string) error {
	idea, e := IdeaSelect(path)
	if e != nil {
		return e
	}
	return cmd.OpenIdea(idea, path)
}
