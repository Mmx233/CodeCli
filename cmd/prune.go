package cmd

import (
	"errors"
	"fmt"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/Mmx233/CodeCli/internal/target"
	"github.com/spf13/cobra"
)

func newPruneCommand(store func() (config.Store, error)) *cobra.Command {
	var apply, force bool
	c := &cobra.Command{Use: "prune <targets...>", Short: "Preview removal of managed checkouts; --apply executes it", Args: cobra.MinimumNArgs(1), RunE: func(c *cobra.Command, args []string) error {
		cfg, err := load(store)
		if err != nil {
			return err
		}
		storage, err := cfg.ProjectDir()
		if err != nil {
			return err
		}
		var paths []string
		for _, raw := range args {
			t, e := target.Parse(raw, cfg)
			if e != nil {
				return e
			}
			if t.Subdir != "" {
				return errors.New("prune expects repository roots, not subdirectories")
			}
			p, e := t.Directory(storage)
			if e != nil {
				return e
			}
			paths = append(paths, p)
		}
		if err = gitFor(c).Prune(c.Context(), storage, paths, apply, force); err != nil {
			return err
		}
		verb := "Would remove"
		if apply {
			verb = "Removed"
		}
		for _, p := range paths {
			if _, err = fmt.Fprintln(c.OutOrStdout(), verb+": "+p); err != nil {
				return err
			}
		}
		if !apply {
			_, err = fmt.Fprintln(c.ErrOrStderr(), "Preview only. Add --apply to remove these checkouts.")
		}
		return err
	}}
	c.Flags().BoolVar(&apply, "apply", false, "Execute the removal shown by preview")
	c.Flags().BoolVar(&force, "force", false, "Allow loss of local files, refs, stash and reflogs; never relax path/worktree boundaries")
	return c
}
