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

// CompleteAddr 填充简写为完整 addr
func CompleteAddr(addr string) (string, error) {
	switch len(strings.Split(addr, "/")) {
	case 1:
		if global.Config.Default.Username == "" {
			return "", util.ErrEmptyDefaultUsername
		}
		addr = global.Config.Default.Username + "/" + addr
		fallthrough
	case 2:
		if global.Config.Default.GitSite == "" {
			return "", util.ErrEmptyDefaultGitSite
		}
		addr = global.Config.Default.GitSite + "/" + addr
	case 3:
		break
	default:
		return "", util.ErrUnknownInput
	}
	return addr, nil
}

func ConvertAddrToPath(addr string) (dir string, path string) {
	split := strings.Split(addr, "/")
	dir = file.JoinPath(append([]string{global.Config.Storage.ProjectDir}, split[:2]...)...)
	path = file.JoinPath(dir, split[2])
	return
}

func Open(addr string) error {
	addr, e := CompleteAddr(addr)
	if e != nil {
		return e
	}

	projectDir, projectPath := ConvertAddrToPath(addr)
	if e = os.MkdirAll(projectDir, 0600); e != nil {
		return e
	}

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
