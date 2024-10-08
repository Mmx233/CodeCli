package browser

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/pkg/project"
	"os/exec"
	"runtime"
)

func Open(addr string) error {
	p, err := project.CompleteAddrToProject(addr)
	if err != nil {
		return err
	}
	return OpenBrowser(p.Url())
}

func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		return exec.Command("open", url).Run()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
