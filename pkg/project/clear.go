package project

import (
	"fmt"
	"github.com/Mmx233/CodeCli/global"
	"os"
	"time"
)

func Clear(t *time.Duration) error {
	entries, e := os.ReadDir(global.Config.Storage.ProjectDir)
	if e != nil {
		return e
	}
	for _, entry := range entries {
		info, e := entry.Info()
		if e != nil {
			return e
		}
		if !info.IsDir() {
			continue
		}
		fmt.Println()
	}
	return nil
}
