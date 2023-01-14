package project

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/tool"
	"os"
	"strings"
)

func Clone(path, url string) error {
	if e := cmd.GitClone(url, path); e != nil {
		_ = os.RemoveAll(path)
		return e
	}
	return nil
}

func CodeUncommitted(path string) (bool, error) {
	r, e := cmd.GitStatus(path)
	if e != nil {
		return false, e
	}
	return !strings.Contains(string(r), "nothing to commit, working tree clean"), nil
}

func LoadProject(addr string) (*Project, error) {
	project, e := CompleteAddrToProject(addr)
	if e != nil {
		return nil, e
	}

	if e = os.MkdirAll(project.Dir, 0600); e != nil {
		return nil, e
	}

	if !tool.File.Exists(project.Path) {
		if e = Clone(project.Path, project.Url()); e != nil {
			return nil, e
		}
	}
	return project, nil
}
