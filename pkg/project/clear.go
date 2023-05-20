package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	"log"
	"os"
	"strings"
	"time"
)

func dirShouldScan(file os.FileInfo) bool {
	return file.IsDir() && !strings.HasPrefix(file.Name(), ".")
}

func Clear(t time.Duration, yes, force bool, addresses ...string) error {
	var projectPaths []string
	var e error
	if len(addresses) != 0 {
		for _, addr := range addresses {
			var project *Project
			project, e = CompleteAddrToProject(addr)
			if e != nil {
				log.Printf("warning: addr %s occur error: %v\n", addr, e)
				continue
			}
			projectPaths = append(projectPaths, project.Path)
		}
	} else {
		//扫描旧项目
		if e = file.ScanDir(global.Config.Storage.ProjectDir, func(path string, info os.FileInfo) error {
			if !dirShouldScan(info) {
				return nil
			}
			return file.ScanDir(file.JoinPath(path, info.Name()), func(path string, info os.FileInfo) error {
				if !dirShouldScan(info) {
					return nil
				}
				return file.ScanDir(file.JoinPath(path, info.Name()), func(path string, info os.FileInfo) error {
					if !dirShouldScan(info) {
						return nil
					}
					path = file.JoinPath(path, info.Name())
					if info.ModTime().Before(time.Now().Add(-t)) {
						projectPaths = append(projectPaths, path)
					}
					return nil
				})
			})
		}); e != nil {
			return e
		}
	}
	if !force && len(projectPaths) != 0 {
		//scan uncommitted repos
		var projectPure []string
		var isClean bool
		for _, path := range projectPaths {
			isClean, e = IsRepoClean(path)
			if e != nil {
				log.Printf("warning: %s isn't a git repo: %v.", path, e)
				continue
			} else if !isClean {
				log.Printf("warning: %s should be cleared, but there are local changes.", path)
			} else {
				projectPure = append(projectPure, path)
			}
		}
		projectPaths = projectPure
	}
	if len(projectPaths) != 0 {
		log.Println("info: following projects is going to be cleared.")
		fmt.Println(strings.Join(projectPaths, "\n"))

		if !yes && !force {
			fmt.Printf("Do you want to continue? [Y/n]")
			var input string
			if _, e = fmt.Scanln(&input); e != nil {
				return e
			}
			if !(input == "" || strings.ToLower(strings.TrimSpace(input)) == "y") {
				return nil
			}
		}

		for _, path := range projectPaths {
			if e = os.RemoveAll(path); e != nil {
				log.Printf("warning: remove project %s failed: %v", path, e)
			}
		}
		log.Println("info: clean task completed.")
	} else {
		log.Println("info: no project to clear.")
	}
	return nil
}
