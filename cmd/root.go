package cmd

import (
	"fmt"
	"os"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	return newRoot(version, config.UserStore)
}

func newRoot(version string, store func() (config.Store, error)) *cobra.Command {
	options := openOptions{}
	root := &cobra.Command{
		Use: "code [target]", Short: "Clone and open projects in an IDE, browser or terminal",
		Version: version, SilenceUsage: true, SilenceErrors: true,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) > 0 {
				switch args[0] {
				case "cmd":
					return fmt.Errorf("code cmd was removed; use code open <target> --terminal")
				case "browser":
					return fmt.Errorf("code browser was removed; use code open <target> --web")
				case "project":
					return fmt.Errorf("code project was removed; use code open <target>")
				case "clear":
					return fmt.Errorf("code clear was removed; use code prune <targets> (preview), then --apply")
				case "path":
					return fmt.Errorf("use code pwd <target>; for a repository named path use code open path")
				case "version":
					return fmt.Errorf("use code --version; for a repository named version use code open version")
				}
			}
			return cobra.MaximumNArgs(1)(c, args)
		},
		RunE: func(c *cobra.Command, args []string) error { return options.run(c, args, store) },
	}
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	root.SetIn(os.Stdin)
	root.Flags().BoolP("version", "v", false, "Print version")
	options.flags(root)
	root.AddCommand(newConfigCommand(store), newOpenCommand(store), newCloneCommand(store), newPwdCommand(store), newPruneCommand(store))
	return root
}
