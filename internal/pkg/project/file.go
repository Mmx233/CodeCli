package project

import (
	"github.com/Mmx233/CodeCli/pkg/git"
	"github.com/Mmx233/tool"
	"os"
	"strings"
)

func Clone(path, url string) error {
	if err := git.Clone(url, path); err != nil {
		_ = os.RemoveAll(path)
		return err
	}
	return nil
}

func IsRepoClean(path string) (bool, error) {
	output, err := git.BranchStatus(path)
	if err != nil {
		return false, err
	}
	if strings.Contains(string(output), "ahead") {
		return false, nil
	}

	output, err = git.Status(path)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), "nothing to commit, working tree clean"), nil
}

func LoadProject(addr string) (*Project, error) {
	project, err := CompleteAddrToProject(addr)
	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(project.Dir, 0775); err != nil {
		return nil, err
	}

	exist, err := tool.File.Exists(project.Path)
	if err != nil {
		return nil, err
	} else if !exist {
		if err = Clone(project.Path, project.Url()); err != nil {
			return nil, err
		}
	}
	return project, nil
}
