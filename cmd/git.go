package cmd

import (
	"os"
	"os/exec"
)

func Clone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if e := cmd.Start(); e != nil {
		return e
	}
	return cmd.Wait()
}
