package project

import (
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/CodeCli/util"
	"strings"
)

// CompleteAddrToUrl 填充简写为完整 addr
func CompleteAddrToUrl(addr string) (string, error) {
	if strings.Contains(addr, "https://") {
		if !strings.HasSuffix(addr, ".git") {
			addr = addr + ".git"
		}
		return addr, nil
	}
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
	return "https://" + addr + ".git", nil
}

func ConvertUrlToPath(addr string) (dir string, path string) {
	addr = strings.TrimLeft(addr, "https://")
	addr = strings.TrimRight(addr, ".git")
	split := strings.Split(addr, "/")
	dir = file.JoinPath(append([]string{global.Config.Storage.ProjectDir}, split[:2]...)...)
	path = file.JoinPath(dir, split[2])
	return
}
