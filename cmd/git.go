package cmd

import (
	"os"
	"os/exec"
)

func Clone(project string) error {
	cmd := exec.Command("git", "clone", project)
	cmd.Stdout = os.Stdout
	if e := cmd.Start(); e != nil {
		return e
	}
	return cmd.Wait()
}
