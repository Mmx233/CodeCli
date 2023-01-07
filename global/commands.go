package global

import (
	"github.com/Mmx233/CodeCli/global/models"
	"gopkg.in/alecthomas/kingpin.v2"
)

var Commands models.Commands

func init() {
	Commands.App = kingpin.New("code", "A project manager command line tool.")
	Commands.App.Version(Version)
	Commands.App.VersionFlag.Short('v')
	Commands.Project.CmdClause = Commands.App.Command("project", "Open projects.").Default()
	Commands.Project.Arg("addr", "Project addr.").Required().HintOptions("github.com/Mmx233/CodeCli").StringVar(&Commands.Project.Addr)

	Commands.Clear.CmdClause = Commands.App.Command("clear", "Auto clear outdated projects.")
	Commands.Clear.Arg("duration", "Clean up projects that have not been used for how long.").Default("1440h").DurationVar(&Commands.Clear.Duration)
	Commands.Clear.Flag("yes", "Confirm delete.").Short('y').BoolVar(&Commands.Clear.Yes)
}

func ParseFlags(args []string) string {
	return kingpin.MustParse(Commands.App.Parse(args))
}
