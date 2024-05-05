package file

import (
	"os"
	"strings"
)

func PreparePath(p string) string {
	p = strings.Replace(p, `\`, `/`, -1)
	return p
}

func ScanDir(path string, f func(path string, info os.FileInfo) error) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		var info os.FileInfo
		info, err = entry.Info()
		if err != nil {
			return err
		}
		if err = f(path, info); err != nil {
			return err
		}
	}
	return nil
}
