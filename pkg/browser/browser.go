package browser

import (
	"github.com/Mmx233/CodeCli/pkg/project"
	"os/exec"
)

func Open(addr string) error {
	p, e := project.CompleteAddrToProject(addr)
	if e != nil {
		return e
	}
	return OpenExplorer(p.Url())
}

func OpenExplorer(path string) error {
	cmd := exec.Command("explorer", path)
	if e := cmd.Start(); e != nil {
		return e
	}
	_ = cmd.Wait()
	return nil
}
