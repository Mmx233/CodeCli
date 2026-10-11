package cmd

import (
	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/Mmx233/CodeCli/internal/project"
	"github.com/Mmx233/CodeCli/internal/target"
	"github.com/spf13/cobra"
)

func load(store func() (config.Store, error)) (config.Config, error) {
	s, err := store()
	if err != nil {
		return config.Config{}, err
	}
	return s.Load()
}

func locate(raw string, store func() (config.Store, error)) (target.Target, string, error) {
	var cfg config.Config
	if !target.IsLocal(raw) {
		var err error
		cfg, err = load(store)
		if err != nil {
			return target.Target{}, "", err
		}
	}
	t, err := target.Parse(raw, cfg)
	if err != nil {
		return t, "", err
	}
	if t.Local {
		return t, "", nil
	}
	storage, err := cfg.ProjectDir()
	return t, storage, err
}

func gitFor(c *cobra.Command) project.Git {
	return project.Git{In: c.InOrStdin(), Err: c.ErrOrStderr()}
}
func argument(args []string) string {
	if len(args) == 0 {
		return "."
	}
	return args[0]
}
