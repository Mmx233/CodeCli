package project

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"github.com/Mmx233/tool"
	"path"
)

func IdeaSelect(dir string) (string, error) {
	for _, rule := range global.Config.Rules {
		for _, filename := range rule.File {
			exist, err := tool.File.Exists(path.Join(dir, filename))
			if err != nil {
				return "", err
			} else if exist {
				return rule.Idea, nil
			}
		}
	}

	if global.Config.Default.Idea != "" {
		return global.Config.Default.Idea, nil
	}
	return "", util.ErrUnsupportedProjectOrEmptyDir{Path: dir}
}
