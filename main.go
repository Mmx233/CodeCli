package main

import (
	"github.com/Mmx233/CodeCli/global"
	"github.com/Mmx233/CodeCli/pkg/config"
	"github.com/Mmx233/CodeCli/pkg/project"
	"log"
	"os"
)

func main() {
	var e error
	switch len(os.Args) {
	case 1:
		e = project.OpenProject(".")
	default:
		switch global.ParseFlags(os.Args[1:]) {
		case global.Commands.Project.FullCommand():
			e = project.Open(global.Commands.Project.Addr)
		case global.Commands.Clear.FullCommand():
			e = project.Clear(global.Commands.Clear.Duration, global.Commands.Clear.Yes)
		case global.Commands.Cmd.FullCommand():
			e = project.OpenCmd(global.Commands.Cmd.Addr)
		case global.Commands.Config.Set.FullCommand():
			e = config.Set(global.Commands.Config.Set.Field, global.Commands.Config.Set.Value)
		case global.Commands.Config.List.FullCommand():
			e = config.List()
		}
	}
	if e != nil {
		log.Fatalln(e)
	}
}
