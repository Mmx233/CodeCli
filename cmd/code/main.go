package main

import (
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/pkg/browser"
	"github.com/Mmx233/CodeCli/pkg/config"
	"github.com/Mmx233/CodeCli/pkg/project"
	"log"
	"os"
)

var Version = "-.-.-"

func init() {
	global.InitCommands(Version)
}

func main() {
	var e error
	switch len(os.Args) {
	case 1:
		var pwd string
		pwd, e = os.Getwd()
		if e == nil {
			e = project.OpenProject(pwd)
		}
	default:
		switch global.ParseFlags(os.Args[1:]) {
		case global.Commands.Project.FullCommand():
			e = project.Open(global.Commands.Project.Addr)
		case global.Commands.Clear.FullCommand():
			e = project.Clear(global.Commands.Clear.Duration,
				global.Commands.Clear.Yes, global.Commands.Clear.Force,
				global.Commands.Clear.Addresses...)
		case global.Commands.Cmd.FullCommand():
			e = project.OpenCmd(global.Commands.Cmd.Addr)
		case global.Commands.Config.List.FullCommand():
			e = config.List()
		case global.Commands.Config.Set.FullCommand():
			e = config.Set(global.Commands.Config.Set.Field, global.Commands.Config.Set.Value)
		case global.Commands.Config.Unset.FullCommand():
			e = config.Unset(global.Commands.Config.Unset.Field)
		case global.Commands.Browser.FullCommand():
			e = browser.Open(global.Commands.Browser.Addr)
		}
	}
	if e != nil {
		log.Fatalln(e)
	}
}
