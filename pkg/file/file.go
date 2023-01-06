package file

import (
	"path"
	"runtime"
	"strings"
)

var JoinPath func(el ...string) string

func init() {
	if runtime.GOOS == "windows" {
		JoinPath = windowsPathJoin
	} else {
		JoinPath = path.Join
	}
}

func windowsPathJoin(el ...string) string {
	return strings.Join(el, `\`)
}
