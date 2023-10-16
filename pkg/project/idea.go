package project

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/tool"
)

const (
	Webstorm      = "webstorm"
	Goland        = "goland"
	AndroidStudio = "studio"
	IntelliJ      = "idea"
	PyCharm       = "pycharm"
)

func IdeaSelect(dir string) (string, error) {
	var selector = map[string][]string{
		Webstorm:      {"package.json"},
		Goland:        {"go.mod"},
		PyCharm:       {"pyproject.toml", "requirements.txt"},
		AndroidStudio: {file.JoinPath("android", "build.gradle"), "build.gradle"},
		IntelliJ:      {"gradlew"},
	}

	for idea, files := range selector {
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
