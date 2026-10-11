package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Mmx233/CodeCli/internal/config"
	"github.com/Mmx233/CodeCli/internal/target"
)

func (g Git) pruneBoundary(ctx context.Context, storage, p string) error {
	if err := target.NoSymlinks(p); err != nil {
		return err
	}
	rel, err := filepath.Rel(storage, p)
	if err != nil {
		return err
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) != 3 || parts[0] == ".." {
		return errors.New("prune only accepts an exact managed host/owner/repo directory")
	}
	host := strings.Replace(parts[0], "~", ":", 1)
	t, err := target.Parse("https://"+host+"/"+parts[1]+"/"+parts[2], config.Config{})
	if err != nil {
		return err
	}
	matched, err := t.Locate(storage)
	if err != nil {
		return err
	}
	if matched != p {
		return errors.New("prune target does not match a unique managed identity")
	}
	if err = g.CheckRoot(ctx, p); err != nil {
		return err
	}
	gitdir := filepath.Join(p, ".git")
	info, err := os.Lstat(gitdir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("prune refuses linked worktrees and external gitdir layouts")
	}
	common, err := g.Output(ctx, p, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return err
	}
	actual, err := filepath.EvalSymlinks(gitdir)
	if err != nil {
		return err
	}
	if filepath.Clean(common) != actual {
		return errors.New("prune refuses shared Git metadata")
	}
	worktrees, err := g.Output(ctx, p, "worktree", "list", "--porcelain")
	if err != nil {
		return err
	}
	if strings.Count("\n"+worktrees, "\nworktree ") != 1 {
		return errors.New("prune refuses repositories with additional worktrees, even with --force")
	}
	return nil
}

func (g Git) pruneData(ctx context.Context, p string) error {
	status, err := g.Output(ctx, p, "status", "--porcelain=v1", "--untracked-files=all", "--ignored=matching")
	if err != nil {
		return err
	}
	if status != "" {
		return errors.New("protected: modified, untracked or ignored files")
	}
	stash, err := g.Output(ctx, p, "for-each-ref", "--format=%(refname)", "refs/stash")
	if err != nil {
		return err
	}
	if stash != "" {
		return errors.New("protected: stash")
	}
	tags, err := g.Output(ctx, p, "for-each-ref", "--format=%(refname)", "refs/tags")
	if err != nil {
		return err
	}
	if tags != "" {
		return errors.New("protected: local tags (their publication cannot be proven without contacting remotes)")
	}
	if _, err = g.Output(ctx, p, "symbolic-ref", "--quiet", "HEAD"); err != nil {
		return errors.New("protected: detached or unborn HEAD")
	}
	commits, err := g.Output(ctx, p, "rev-list", "--all", "--reflog", "--not", "--remotes")
	if err != nil {
		return err
	}
	if commits != "" {
		return errors.New("protected: commits reachable from local refs or reflogs but absent from remote-tracking refs")
	}
	remotes, err := g.Output(ctx, p, "for-each-ref", "--format=%(refname)", "refs/remotes")
	if err != nil {
		return err
	}
	names := map[string]bool{}
	for ref := range strings.SplitSeq(remotes, "\n") {
		_, name, ok := strings.Cut(strings.TrimPrefix(ref, "refs/remotes/"), "/")
		if ok {
			names[name] = true
		}
	}
	branches, err := g.Output(ctx, p, "for-each-ref", "--format=%(refname:strip=2)", "refs/heads")
	if err != nil {
		return err
	}
	for branch := range strings.SplitSeq(branches, "\n") {
		if !names[branch] {
			return errors.New("protected: local branch name absent from remote-tracking branches")
		}
	}
	refs, err := g.Output(ctx, p, "for-each-ref", "--format=%(refname)")
	if err != nil {
		return err
	}
	for ref := range strings.SplitSeq(refs, "\n") {
		if !strings.HasPrefix(ref, "refs/heads/") && !strings.HasPrefix(ref, "refs/remotes/") && ref != "" {
			return errors.New("protected: additional local Git refs")
		}
	}
	return nil
}

// Prune validates the complete batch before touching any checkout.
func (g Git) Prune(ctx context.Context, storage string, paths []string, apply, force bool) error {
	if err := target.NoSymlinks(storage); err != nil {
		return err
	}
	var unlock func()
	if apply {
		var err error
		unlock, err = lock(storage)
		if err != nil {
			return err
		}
		defer unlock()
	}
	seen := map[string]bool{}
	for _, p := range paths {
		if seen[p] {
			return errors.New("duplicate prune target")
		}
		seen[p] = true
		if err := g.pruneBoundary(ctx, storage, p); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		if !force {
			if err := g.pruneData(ctx, p); err != nil {
				return fmt.Errorf("%s: %w; inspect the data before using --force", p, err)
			}
		}
	}
	if !apply {
		return nil
	}
	root, err := os.OpenRoot(storage)
	if err != nil {
		return err
	}
	defer root.Close()
	for _, p := range paths {
		if err = ctx.Err(); err != nil {
			return err
		}
		// Recheck immediately before removal; rooted operations cannot escape storage.
		if err = g.pruneBoundary(ctx, storage, p); err != nil {
			return err
		}
		if !force {
			if err = g.pruneData(ctx, p); err != nil {
				return err
			}
		}
		rel, err := filepath.Rel(storage, p)
		if err != nil {
			return err
		}
		parent, err := root.OpenRoot(filepath.Dir(rel))
		if err != nil {
			return err
		}
		err = parent.RemoveAll(filepath.Base(rel))
		parent.Close()
		if err != nil {
			return fmt.Errorf("remove %s: %w", p, err)
		}
	}
	return nil
}
