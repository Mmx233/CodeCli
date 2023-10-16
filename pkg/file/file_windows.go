package file

import "strings"

func init() {
	JoinPath = windowsPathJoin
}

func windowsPathJoin(el ...string) string {
	return strings.Join(el, `\`)
}
