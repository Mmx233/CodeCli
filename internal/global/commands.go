package global

import (
	"github.com/Mmx233/CodeCli/internal/global/models"
	"github.com/alecthomas/kingpin/v2"
)

var Commands models.Commands

func InitCommands(Version string) {
	Commands.App = kingpin.New("code", "A project manager command line tool.")
	Commands.App.Version(Version)
	Commands.App.VersionFlag.Short('v')
	Commands.App.HelpFlag.Short('h')

	Commands.Project.CmdClause = Commands.App.Command("project", "Open projects.").Default()
	Commands.Project.Arg("addr", "Project addr.").Default(".").HintOptions("github.com/Mmx233/CodeCli").StringVar(&Commands.Project.Addr)
	Commands.Project.Flag("idea", "Specify an idea.").HintOptions("goland", "webstorm").StringVar(&Commands.Project.Idea)

	Commands.Clear.CmdClause = Commands.App.Command("clear", "Auto clear outdated projects.")
	Commands.Clear.Flag("time", "Clean up projects that have not been used for how long.").Short('t').Default("1440h").DurationVar(&Commands.Clear.Duration)
	Commands.Clear.Flag("yes", "Confirm delete.").Short('y').BoolVar(&Commands.Clear.Yes)
	Commands.Clear.Flag("force", "Force delete, skip confirm and checks.").Short('f').BoolVar(&Commands.Clear.Force)
	Commands.Clear.Arg("addr", "Project addr.").StringsVar(&Commands.Clear.Addresses)

	Commands.Cmd.CmdClause = Commands.App.Command("cmd", "Open project terminal.")
	Commands.Cmd.Arg("addr", "Project addr.").Default(".").StringVar(&Commands.Cmd.Addr)

	Commands.Config.CmdClause = Commands.App.Command("config", "Write configs.")
	Commands.Config.List.CmdClause = Commands.Config.Command("list", "List all configs.")
	Commands.Config.Set.CmdClause = Commands.Config.Command("set", "Set config.").Default()
	Commands.Config.Set.Arg("field", "Field of config.").Required().StringVar(&Commands.Config.Set.Field)
	Commands.Config.Set.Arg("value", "Value of field.").StringVar(&Commands.Config.Set.Value)
	Commands.Config.Unset.CmdClause = Commands.Config.Command("unset", "Clear config.")
	Commands.Config.Unset.Arg("field", "Field to clear.").Required().StringVar(&Commands.Config.Unset.Field)

	Commands.Browser.CmdClause = Commands.App.Command("browser", "Open project in browser.")
	Commands.Browser.Arg("addr", "Project addr.").Default(".").StringVar(&Commands.Browser.Addr)
}

func ParseFlags(args []string) string {
	return kingpin.MustParse(Commands.App.Parse(args))
}
