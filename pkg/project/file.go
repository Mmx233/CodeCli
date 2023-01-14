package project

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/tool"
	"os"
)

func Clone(path, url string) error {
	if e := cmd.Clone(url, path); e != nil {
		_ = os.RemoveAll(path)
		return e
	}
	return nil
}

func LoadProject(addr string) (url string, dir string, path string, e error) {
	url, e = CompleteAddrToUrl(addr)
	if e != nil {
		return
	}

	dir, path = ConvertUrlToPath(url)
	if e = os.MkdirAll(dir, 0600); e != nil {
		return
	}

	if !tool.File.Exists(path) {
		if e = Clone(path, url); e != nil {
			return
		}
	}
	return
}
