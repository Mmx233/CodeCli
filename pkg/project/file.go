package project

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/tool"
	"os"
)

func Clone(path, addr string) error {
	if !tool.File.Exists(path) {
		return cmd.Clone("https://"+addr+".git", path)
	}
	return nil
}

func PrepareProjectFiles(addr string) (fullAddr string, dir string, path string, e error) {
	fullAddr, e = CompleteAddr(addr)
	if e != nil {
		return
	}

	dir, path = ConvertAddrToPath(fullAddr)
	if e = os.MkdirAll(dir, 0600); e != nil {
		return
	}

	if e = Clone(path, fullAddr); e != nil {
		return
	}
	return
}
