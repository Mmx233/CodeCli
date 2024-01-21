package project

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/tool"
)

func IdeaSelect(dir string) (string, error) {
	for idea, files := range global.Config.Rules {
		for _, filename := range files {
			exist, err := tool.File.Exists(file.JoinPath(dir, filename))
			if err != nil {
				return "", err
			} else if exist {
				return idea, nil
			}
		}
	}

	if global.Config.Default.Idea != "" {
		return global.Config.Default.Idea, nil
	}
	return "", util.ErrUnsupportedProjectOrEmptyDir{Path: dir}
}
