package cmd

import (
	"fmt"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCommand(store func() (config.Store, error)) *cobra.Command {
	c := &cobra.Command{Use: "config", Short: "Read or change ~/.CodeCli.yaml"}
	c.AddCommand(&cobra.Command{
		Use: "get [key]", Short: "Print effective configuration (or one scalar key)", Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			s, err := store()
			if err != nil {
				return err
			}
			cfg, err := s.Load()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				v, err := cfg.Value(args[0])
				if err != nil {
					return err
				}
				_, err = fmt.Fprintln(c.OutOrStdout(), v)
				return err
			}
			data, err := cfg.YAML()
			if err != nil {
				return err
			}
			_, err = c.OutOrStdout().Write(data)
			return err
		},
	})
	c.AddCommand(&cobra.Command{
		Use: "set <key> <value>", Short: "Set a scalar key; rules are edited manually", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			s, err := store()
			if err != nil {
				return err
			}
			return s.Set(args[0], &args[1])
		},
	})
	c.AddCommand(&cobra.Command{
		Use: "unset <key>", Short: "Remove an override and restore its default", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			s, err := store()
			if err != nil {
				return err
			}
			return s.Set(args[0], nil)
		},
	})
	return c
}
