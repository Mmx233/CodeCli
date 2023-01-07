package global

import (
	"github.com/Mmx233/CodeCli/global/models"
	"gopkg.in/alecthomas/kingpin.v2"
)

var Commands models.Commands

func init() {
	kingpin.Version(Version)
	Commands.App = kingpin.New("code", "A project manager command line tool.")
	Commands.Project.CmdClause = Commands.App.Command("project", "Open projects.").Default()
	Commands.Project.Arg("addr", "Project addr.").Required().HintOptions("github.com/Mmx233/CodeCli").StringVar(&Commands.Project.Addr)

	Commands.Clear.CmdClause = Commands.App.Command("clear", "Auto clear outdated projects.")
	Commands.Clear.Duration = Commands.Clear.Arg("duration", "Clean up projects that have not been used for how long.").Default("1440h").Duration()
}

func ParseFlags(args []string) string {
	return kingpin.MustParse(Commands.App.Parse(args))
}
