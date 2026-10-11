package cmd

import (
	"fmt"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/spf13/cobra"
)

func newCloneCommand(store func() (config.Store, error)) *cobra.Command {
	var fetch, pull, rebase bool
	c := &cobra.Command{Use: "clone <target>", Short: "Clone or reuse a checkout; optionally update an existing one", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		t, storage, err := locate(args[0], store)
		if err != nil {
			return err
		}
		mode := ""
		if fetch {
			mode = "fetch"
		}
		if pull {
			mode = "pull"
		}
		if rebase {
			mode = "rebase"
		}
		p, err := gitFor(c).Clone(c.Context(), t, storage, mode)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(c.OutOrStdout(), p)
		return err
	}}
	c.Flags().BoolVar(&fetch, "fetch", false, "Fetch existing checkout without changing the working tree")
	c.Flags().BoolVar(&pull, "pull", false, "Update existing checkout with a fast-forward-only pull")
	c.Flags().BoolVar(&rebase, "rebase", false, "Update existing checkout with pull --rebase")
	c.MarkFlagsMutuallyExclusive("fetch", "pull", "rebase")
	return c
}
