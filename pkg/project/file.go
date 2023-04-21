package project

import (
	"github.com/Mmx233/CodeCli/pkg/git"
	"github.com/Mmx233/tool"
	"os"
	"strings"
)

func Clone(path, url string) error {
	if e := git.Clone(url, path); e != nil {
		_ = os.RemoveAll(path)
		return e
	}
	return nil
}

func CodeUncommitted(path string) (bool, error) {
	r, e := git.Status(path)
	if e != nil {
		return false, e
	}
	str := string(r)

	return !(strings.Contains(str, "nothing to commit, working tree clean") &&
			strings.Contains(str, "Your branch is up to date with")),
		nil
}

func LoadProject(addr string) (*Project, error) {
	project, e := CompleteAddrToProject(addr)
	if e != nil {
		return nil, e
	}

	if e = os.MkdirAll(project.Dir, 0600); e != nil {
		return nil, e
	}

	exist, e := tool.File.Exists(project.Path)
	if e != nil {
		return nil, e
	} else if !exist {
		if e = Clone(project.Path, project.Url()); e != nil {
			return nil, e
		}
	}
	return project, nil
}
