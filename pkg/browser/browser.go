package browser

import (
	"github.com/Mmx233/CodeCli/cmd"
	"github.com/Mmx233/CodeCli/pkg/project"
)

func Open(addr string) error {
	p, e := project.CompleteAddrToProject(addr)
	if e != nil {
		return e
	}
	return cmd.OpenExplorer(p.Url())
}
