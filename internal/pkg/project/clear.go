package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/file"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

func dirShouldScan(file os.FileInfo) bool {
	return file.IsDir() && !strings.HasPrefix(file.Name(), ".")
}

func Clear(t time.Duration, yes, force bool, addresses ...string) error {
	logger := log.WithField("component", "clear")

	var projectPaths []string
	var err error
	if len(addresses) != 0 {
		for _, addr := range addresses {
			var project *Project
			project, err = CompleteAddrToProject(addr)
			if err != nil {
				logger.Warnf("addr %s occur error: %v", addr, err)
				continue
			}
			projectPaths = append(projectPaths, project.Path)
		}
	} else {
		//扫描旧项目
		if err = file.ScanDir(global.Config.Storage.ProjectDir, func(dir string, info os.FileInfo) error {
			if !dirShouldScan(info) {
				return nil
			}
			return file.ScanDir(path.Join(dir, info.Name()), func(dir string, info os.FileInfo) error {
				if !dirShouldScan(info) {
					return nil
				}
				return file.ScanDir(path.Join(dir, info.Name()), func(dir string, info os.FileInfo) error {
					if !dirShouldScan(info) {
						return nil
					}
					dir = path.Join(dir, info.Name())
					if info.ModTime().Before(time.Now().Add(-t)) {
						projectPaths = append(projectPaths, dir)
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
		for _, projectPath := range projectPaths {
			isClean, err = IsRepoClean(projectPath)
			if err != nil {
				logger.Warnf("%s isn't a git repo: %v.", projectPath, err)
				continue
			} else if !isClean {
				logger.Warnf("%s should be cleared, but there are local changes.", projectPath)
			} else {
				projectPure = append(projectPure, projectPath)
			}
		}
		projectPaths = projectPure
	}
	if len(projectPaths) != 0 {
		logger.Infoln("following projects is going to be cleared.")
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

		for _, projectPath := range projectPaths {
			if err = os.RemoveAll(projectPath); err != nil {
				logger.Warnf("remove project %s failed: %v", projectPath, err)
			}
		}
		logger.Infoln("clean task completed.")
	} else {
		logger.Infoln("no project to clear.")
	}
	return nil
}
