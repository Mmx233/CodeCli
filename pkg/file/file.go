package file

import (
	"os"
)

var JoinPath func(el ...string) string

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
