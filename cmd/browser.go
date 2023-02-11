package cmd

import (
	"os/exec"
)

func OpenExplorer(path string) error {
	cmd := exec.Command("explorer", path)
	if e := cmd.Start(); e != nil {
		return e
	}
	_ = cmd.Wait()
	return nil
}
