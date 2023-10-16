package file

import (
	"os"
)

var JoinPath func(el ...string) string

func ScanDir(path string, f func(path string, info os.FileInfo) error) error {
	entries, e := os.ReadDir(path)
	if e != nil {
		return e
	}
	for _, entry := range entries {
		var info os.FileInfo
		info, e = entry.Info()
		if e != nil {
			return e
		}
		if e = f(path, info); e != nil {
			return e
		}
	}
	return nil
}
