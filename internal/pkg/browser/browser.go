package browser

import (
	"github.com/Mmx233/CodeCli/internal/pkg/project"
	"os/exec"
)

func Open(addr string) error {
	p, err := project.CompleteAddrToProject(addr)
	if err != nil {
		return err
	}
	return OpenExplorer(p.Url())
}

func OpenExplorer(path string) error {
	cmd := exec.Command("explorer", path)
	if err := cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Wait()
	return nil
}
