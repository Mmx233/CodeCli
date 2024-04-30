package main

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/pkg/browser"
	"github.com/Mmx233/CodeCli/internal/pkg/config"
	"github.com/Mmx233/CodeCli/internal/pkg/project"
	"log"
	"os"
)

var Version = "-.-.-"

func init() {
	global.InitCommands(Version)
}

func main() {
	var err error
	switch len(os.Args) {
	case 1:
		var pwd string
		pwd, err = os.Getwd()
		if err == nil {
			err = project.OpenProject(pwd)
		}
	default:
		switch global.ParseFlags(os.Args[1:]) {
		case global.Commands.Project.FullCommand():
			err = project.Open(global.Commands.Project.Addr)
		case global.Commands.Clear.FullCommand():
			err = project.Clear(global.Commands.Clear.Duration,
				global.Commands.Clear.Yes, global.Commands.Clear.Force,
				global.Commands.Clear.Addresses...)
		case global.Commands.Cmd.FullCommand():
			err = project.OpenCmd(global.Commands.Cmd.Addr)
		case global.Commands.Config.List.FullCommand():
			err = config.List()
		case global.Commands.Config.Set.FullCommand():
			err = config.Set(global.Commands.Config.Set.Field, global.Commands.Config.Set.Value)
		case global.Commands.Config.Unset.FullCommand():
			err = config.Unset(global.Commands.Config.Unset.Field)
		case global.Commands.Browser.FullCommand():
			err = browser.Open(global.Commands.Browser.Addr)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
}
