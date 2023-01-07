package project

import (
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/CodeCli/util"
	"github.com/Mmx233/tool"
)

const (
	Webstorm      = "webstorm"
	Goland        = "goland"
	AndroidStudio = "studio"
)

func IdeaSelect(dir string) (string, error) {
	if tool.File.Exists(file.JoinPath(dir, "package.json")) {
		return Webstorm, nil
	}
	if tool.File.Exists(file.JoinPath(dir, "go.mod")) {
		return Goland, nil
	}
	if tool.File.Exists(file.JoinPath(dir, "android", "build.gradle")) || tool.File.Exists(file.JoinPath(dir, "build.gradle")) {
		return AndroidStudio, nil
	}
	if global.Config.Default.Idea != "" {
		return global.Config.Default.Idea, nil
	}
	return "", util.ErrUnsupportedProjectOrEmptyDir{Path: dir}
}
