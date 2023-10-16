//go:build !windows

package file

import "path"

func init() {
	JoinPath = path.Join
}
