package cmd

import "os/exec"

func OpenExplorer(path string) error {
	return exec.Command("explorer", path).Run()
}
