package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Clear(t time.Duration, yes bool) error {
	var outdatedProjects []string
	e := file.ScanDir(global.Config.Storage.ProjectDir, func(path string, info os.FileInfo) error {
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
				if info.ModTime().Before(time.Now().Add(-t)) {
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
						outdatedProjects = append(outdatedProjects, path)
						return nil
					}
				}
				return nil
			})
		})
	})
	if e != nil {
		return e
	}
	if len(outdatedProjects) != 0 {
		log.Println("info: following projects is going to be deleted.")
		fmt.Println(strings.Join(outdatedProjects, "\n"))

		if !yes {
			fmt.Printf("Do you want to continue? [Y/n]")
			var input string
			if _, e = fmt.Scanln(&input); e != nil {
				return e
			}
			if !(input == "" || strings.ToLower(strings.TrimSpace(input)) == "y") {
				return nil
			}
		}

		for _, path := range outdatedProjects {
			if e = os.RemoveAll(path); e != nil {
				log.Printf("warning: remove project %s failed: %v", path, e)
			}
		}
		log.Println("info: clear task completed.")
	} else {
		log.Println("info: no project is outdated.")
	}
	return nil
}
