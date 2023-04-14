package git

import (
	"fmt"
	"os"
	"os/exec"
)

func Clone(url, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Cloning %s\n", url)
	if e := cmd.Start(); e != nil {
		return e
	}
	return cmd.Wait()
}

func Status(path string) ([]byte, error) {
	cmd := exec.Command("git", "status")
	cmd.Dir = path
	return cmd.Output()
}
