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
)

func IdeaSelect(dir string) (string, error) {
	exist, e := tool.File.Exists(file.JoinPath(dir, "package.json"))
	if e != nil {
		return "", e
	} else if exist {
		return Webstorm, nil
	}

	exist, e = tool.File.Exists(file.JoinPath(dir, "go.mod"))
	if e != nil {
		return "", e
	} else if exist {
		return Goland, nil
	}

	exist, e = tool.File.Exists(file.JoinPath(dir, "android", "build.gradle"))
	if e != nil {
		return "", e
	} else if exist {
		return AndroidStudio, nil
	} else if exist, e = tool.File.Exists(file.JoinPath(dir, "build.gradle")); e != nil {
		return "", e
	} else if exist {
		return AndroidStudio, nil
	}

	exist, e = tool.File.Exists(file.JoinPath(dir, "gradlew"))
	if e != nil {
		return "", e
	} else if exist {
		return IntelliJ, nil
	}

	if global.Config.Default.Idea != "" {
		return global.Config.Default.Idea, nil
	}
	return "", util.ErrUnsupportedProjectOrEmptyDir{Path: dir}
}
