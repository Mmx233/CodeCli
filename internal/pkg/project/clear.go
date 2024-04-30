package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func dirShouldScan(file os.FileInfo) bool {
	return file.IsDir() && !strings.HasPrefix(file.Name(), ".")
}

func Clear(t time.Duration, yes, force bool, addresses ...string) error {
	var projectPaths []string
	var err error
	if len(addresses) != 0 {
		for _, addr := range addresses {
			var project *Project
			project, err = CompleteAddrToProject(addr)
			if err != nil {
				log.Warnf("warning: addr %s occur error: %v\n", addr, err)
				continue
			}
			projectPaths = append(projectPaths, project.Path)
		}
	} else {
		//扫描旧项目
		if err = file.ScanDir(global.Config.Storage.ProjectDir, func(path string, info os.FileInfo) error {
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
		}); err != nil {
			return err
		}
	}
	if !force && len(projectPaths) != 0 {
		//scan uncommitted repos
		var projectPure []string
		var isClean bool
		for _, path := range projectPaths {
			isClean, err = IsRepoClean(path)
			if err != nil {
				log.Warnf("%s isn't a git repo: %v.", path, err)
				continue
			} else if !isClean {
				log.Warnf("%s should be cleared, but there are local changes.", path)
			} else {
				projectPure = append(projectPure, path)
			}
		}
		projectPaths = projectPure
	}
	if len(projectPaths) != 0 {
		log.Infoln("following projects is going to be cleared.")
		fmt.Println(strings.Join(projectPaths, "\n"))

		if !yes && !force {
			fmt.Printf("Do you want to continue? [Y/n]")
			var input string
			if _, err = fmt.Scanln(&input); err != nil {
				return err
			}
			if !(input == "" || strings.ToLower(strings.TrimSpace(input)) == "y") {
				return nil
			}
		}

		for _, path := range projectPaths {
			if err = os.RemoveAll(path); err != nil {
				log.Printf("warning: remove project %s failed: %v", path, err)
			}
		}
		log.Infoln("clean task completed.")
	} else {
		log.Infoln("no project to clear.")
	}
	return nil
}
