package project

import (
	"github.com/Mmx233/CodeCli/util"
	"github.com/Mmx233/tool"
	"path"
)

const (
	Webstorm = "webstorm"
	Goland   = "goland"
)

func IdeaSelect(dir string) (string, error) {
	if tool.File.Exists(path.Join(dir, "package.json")) {
		return Webstorm, nil
	}
	if tool.File.Exists(path.Join(dir, "go.mod")) {
		return Goland, nil
	}
	return "", util.ErrEmptyDir{Path: dir}
}
