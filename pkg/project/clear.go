package project

import (
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Clear(t *time.Duration) error {
	return file.ScanDir(global.Config.Storage.ProjectDir, func(path string, info os.FileInfo) error {
		if !info.IsDir() {
			return nil
		}
		return file.ScanDir(file.JoinPath(path, info.Name()), func(path string, info os.FileInfo) error {
			if !info.IsDir() {
				return nil
			}
			return file.ScanDir(file.JoinPath(path, info.Name()), func(path string, info os.FileInfo) error {
				if !info.IsDir() {
					return nil
				}
				path = file.JoinPath(path, info.Name())
				if info.ModTime().Before(time.Now().Add(-*t)) {
					cmd := exec.Command("git", "status")
					cmd.Dir = path
					r, e := cmd.Output()
					if e != nil {
						log.Printf("warning: %s isn't a git repo: %v.", path, e)
						return nil
					}
					if !strings.Contains(string(r), "nothing to commit, working tree clean") {
						log.Printf("warning: %s is outdated but have uncommited codes.", path)
						return nil
					} else {
						log.Printf("info: cleaing %s", path)
						return os.RemoveAll(path)
					}
				}
				return nil
			})
		})
	})
}
