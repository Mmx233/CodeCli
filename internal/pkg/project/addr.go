package project

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"github.com/Mmx233/CodeCli/pkg/file"
	"strings"
)

// CompleteAddrToProject 填充简写为完整 addr
func CompleteAddrToProject(addr string) (*Project, error) {
	addr = strings.Replace(addr, `\`, `/`, -1)
	if strings.Contains(addr, "https://") {
		addr = strings.TrimPrefix(addr, "https://")
		addr = strings.TrimSuffix(addr, ".git")
	}
	var p Project
	infos := strings.Split(addr, "/")
	switch len(infos) {
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
	case 3:
		p.GitSite = infos[0]
		p.Username = infos[1]
		p.Repo = infos[2]
	default:
		return nil, util.ErrIllegalInput
	}

	p.Dir = file.JoinPath(global.Config.Storage.ProjectDir, p.GitSite, p.Username)
	p.Path = file.JoinPath(p.Dir, p.Repo)
	return &p, nil
}
