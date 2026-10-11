package cmd

import (
	"fmt"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/spf13/cobra"
)

func newPwdCommand(store func() (config.Store, error)) *cobra.Command {
	return &cobra.Command{Use: "pwd [target]", Short: "Print an existing absolute directory; never clone or update", Args: cobra.MaximumNArgs(1), RunE: func(c *cobra.Command, args []string) error {
		t, storage, err := locate(argument(args), store)
		if err != nil {
			return err
		}
		p, err := t.Directory(storage)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(c.OutOrStdout(), p)
		return err
	}}
}
