package project

import (
	"errors"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"github.com/Mmx233/CodeCli/pkg/file"
	"os"
	"path"
	"strings"
)

// CompleteAddrToProject 填充简写为完整 addr
func CompleteAddrToProject(addr string) (*Project, error) {
	if addr == "" {
		addr = "."
	} else {
		addr = file.PreparePath(addr)
	}

	switch {
	case strings.HasPrefix(addr, "http://"):
		addr = strings.TrimPrefix(addr, "http://")
		addr = strings.TrimSuffix(addr, ".git")
	case strings.HasPrefix(addr, "https://"):
		addr = strings.TrimPrefix(addr, "https://")
		addr = strings.TrimSuffix(addr, ".git")
	case addr == "." || addr == ".." ||
		strings.HasPrefix(addr, "./") || strings.HasPrefix(addr, "../"):
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		addr = file.PreparePath(path.Join(pwd, addr))
		fallthrough
	case path.IsAbs(addr):
		addr = path.Clean(addr)
		if !strings.HasPrefix(addr, global.Config.Storage.ProjectDir) {
			return nil, errors.New("target path out of project root")
		}
		addr = strings.TrimPrefix(addr, global.Config.Storage.ProjectDir)
	default:
		return nil, util.ErrIllegalInput
	}

	var p Project
	infos := strings.Split(addr, "/")
	switch len(infos) {
	case 0:
		return nil, util.ErrIllegalInput
	case 1:
		if global.Config.Default.Username == "" {
			return nil, util.ErrEmptyDefaultUsername
		} else if global.Config.Default.GitSite == "" {
			return nil, util.ErrEmptyDefaultGitSite
		}
		p.GitSite = global.Config.Default.GitSite
		p.Username = global.Config.Default.Username
		p.Repo = infos[0]
	case 2:
		if global.Config.Default.GitSite == "" {
			return nil, util.ErrEmptyDefaultGitSite
		}
		p.GitSite = global.Config.Default.GitSite
		p.Username = infos[0]
		p.Repo = infos[1]
	default:
		p.SubDir = path.Join(infos[3:]...)
		fallthrough
	case 3:
		p.GitSite = infos[0]
		p.Username = infos[1]
		p.Repo = infos[2]
	}

	p.Dir = path.Join(global.Config.Storage.ProjectDir, p.GitSite, p.Username)
	p.Path = path.Join(p.Dir, p.Repo, p.SubDir)
	pState, err := os.Stat(p.Path)
	if err != nil {
		return nil, err
	}
	if !pState.IsDir() {
		return nil, util.ErrProjectMustDir
	}
	return &p, nil
}
