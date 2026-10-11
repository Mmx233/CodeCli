package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/Mmx233/CodeCli/internal/desktop"
	"github.com/Mmx233/CodeCli/internal/target"
	"github.com/spf13/cobra"
)

type openOptions struct {
	web, terminal bool
	ide           string
}

func (o *openOptions) flags(c *cobra.Command) {
	c.Flags().BoolVar(&o.web, "web", false, "Open the repository homepage without cloning")
	c.Flags().BoolVar(&o.terminal, "terminal", false, "Open a desktop terminal in the project")
	c.Flags().StringVar(&o.ide, "ide", "", "IDE executable or macOS application path")
	c.MarkFlagsMutuallyExclusive("web", "terminal", "ide")
}

func newOpenCommand(store func() (config.Store, error)) *cobra.Command {
	o := openOptions{}
	c := &cobra.Command{Use: "open [target]", Short: "Open a project (default: current directory)", Args: cobra.MaximumNArgs(1), RunE: func(c *cobra.Command, args []string) error { return o.run(c, args, store) }}
	o.flags(c)
	return c
}

func (o openOptions) run(c *cobra.Command, args []string, store func() (config.Store, error)) error {
	raw := argument(args)
	t, storage, err := locate(raw, store)
	if err != nil {
		return err
	}
	g := gitFor(c)
	if o.web {
		if t.Local {
			if _, err = target.ExistingDir(t.Path); err != nil {
				return err
			}
			u, e := g.Origin(c.Context(), t.Path)
			if e != nil {
				return e
			}
			t, err = target.Parse(u, config.Config{})
			if err != nil {
				return fmt.Errorf("origin cannot be mapped to a webpage: %w", err)
			}
		}
		address, e := t.WebURL()
		if e != nil {
			return e
		}
		if err = desktop.CheckSession(); err != nil {
			fmt.Fprintln(c.ErrOrStderr(), address)
			return err
		}
		p, e := desktop.Web(runtime.GOOS, address, exec.LookPath)
		if e != nil {
			fmt.Fprintln(c.ErrOrStderr(), address)
			return e
		}
		return p.Launch(c.Context(), c.ErrOrStderr())
	}
	if err = desktop.CheckSession(); err != nil {
		return err
	}
	cfg, err := load(store)
	if err != nil {
		return err
	}
	var plan desktop.Plan
	// Known launchers are checked before any clone; automatic rules need the checkout.
	if o.terminal {
		plan, err = desktop.Terminal(runtime.GOOS, cfg.Default.CmdProgram, "", exec.LookPath)
	} else if o.ide != "" {
		plan, err = desktop.IDE(runtime.GOOS, o.ide, "", exec.LookPath)
	}
	if err != nil {
		return err
	}
	dir, err := g.Ensure(c.Context(), t, storage)
	if err != nil {
		return err
	}
	if o.terminal {
		plan, err = desktop.Terminal(runtime.GOOS, cfg.Default.CmdProgram, dir, exec.LookPath)
	} else {
		name, e := desktop.SelectIDE(cfg, o.ide, dir)
		if e != nil {
			return e
		}
		plan, err = desktop.IDE(runtime.GOOS, name, dir, exec.LookPath)
	}
	if err != nil {
		return err
	}
	return plan.Launch(c.Context(), c.ErrOrStderr())
}
