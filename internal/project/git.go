package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Mmx233/CodeCli/internal/target"
)

type Git struct {
	In  io.Reader
	Err io.Writer
}

func (g Git) command(ctx context.Context, dir string, args ...string) *exec.Cmd {
	c := exec.CommandContext(ctx, "git", args...)
	c.Dir = dir
	c.Stdin = g.In
	// A caller's GIT_DIR must not redirect safety checks to a different checkout.
	for _, e := range os.Environ() {
		key, _, _ := strings.Cut(e, "=")
		switch key {
		case "GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES":
			continue
		}
		c.Env = append(c.Env, e)
	}
	return c
}

func (g Git) Output(ctx context.Context, dir string, args ...string) (string, error) {
	c := g.command(ctx, dir, args...)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	b, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", args[0], err, target.RedactText(strings.TrimSpace(stderr.String())))
	}
	return strings.TrimSpace(string(b)), nil
}

func (g Git) Run(ctx context.Context, dir string, args ...string) error {
	c := g.command(ctx, dir, args...)
	writer := g.Err
	if writer == nil {
		writer = io.Discard
	}
	out := &redactedWriter{dst: writer}
	c.Stdout = out
	c.Stderr = out
	err := c.Run()
	if e := out.flush(); err == nil {
		err = e
	}
	if err != nil {
		return fmt.Errorf("git %s failed: %w", args[0], err)
	}
	return nil
}

// Keep only an unfinished word so split writes cannot expose URL credentials.
type redactedWriter struct {
	dst     io.Writer
	pending string
}

func (w *redactedWriter) Write(p []byte) (int, error) {
	w.pending += string(p)
	end := strings.LastIndexAny(w.pending, " \t\r\n'\"<>")
	if end < 0 {
		return len(p), nil
	}
	if _, err := io.WriteString(w.dst, target.RedactText(w.pending[:end+1])); err != nil {
		return 0, err
	}
	w.pending = w.pending[end+1:]
	return len(p), nil
}

func (w *redactedWriter) flush() error {
	_, err := io.WriteString(w.dst, target.RedactText(w.pending))
	w.pending = ""
	return err
}

func (g Git) Root(ctx context.Context, dir string) (string, error) {
	p, err := g.Output(ctx, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", errors.New("target must be a Git checkout")
	}
	return filepath.Abs(p)
}

func (g Git) CheckRoot(ctx context.Context, dir string) error {
	root, err := g.Root(ctx, dir)
	if err != nil {
		return err
	}
	a, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}
	b, err := filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	if a != b {
		return errors.New("target must be the Git checkout root, not a subdirectory")
	}
	return nil
}

func (g Git) Origin(ctx context.Context, dir string) (string, error) {
	root, err := g.Root(ctx, dir)
	if err != nil {
		return "", err
	}
	u, err := g.Output(ctx, root, "remote", "get-url", "origin")
	if err != nil {
		return "", errors.New("checkout has no origin; provide an explicit HTTPS target")
	}
	return u, nil
}

func (g Git) Update(ctx context.Context, dir, mode string) error {
	if err := g.CheckRoot(ctx, dir); err != nil {
		return err
	}
	if mode == "" {
		return nil
	}
	if mode == "fetch" {
		return g.Run(ctx, dir, "fetch")
	}
	if mode != "pull" && mode != "rebase" {
		return errors.New("unknown update mode")
	}
	if _, err := g.Output(ctx, dir, "symbolic-ref", "--quiet", "HEAD"); err != nil {
		return errors.New("cannot update a detached HEAD; switch to a branch")
	}
	if _, err := g.Output(ctx, dir, "rev-parse", "--verify", "@{upstream}"); err != nil {
		return errors.New("current branch has no upstream; configure one before updating")
	}
	status, err := g.Output(ctx, dir, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return err
	}
	if status != "" {
		return errors.New("working tree has changes or untracked files; commit or move them before updating")
	}
	for _, state := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "REBASE_HEAD", "rebase-merge", "rebase-apply", "sequencer", "BISECT_LOG"} {
		p, e := g.Output(ctx, dir, "rev-parse", "--git-path", state)
		if e != nil {
			return e
		}
		if !filepath.IsAbs(p) {
			p = filepath.Join(dir, p)
		}
		if _, e = os.Lstat(p); e == nil {
			return errors.New("Git operation is in progress; finish or abort it before updating")
		} else if !errors.Is(e, os.ErrNotExist) {
			return e
		}
	}
	if mode == "pull" {
		return g.Run(ctx, dir, "-c", "pull.rebase=false", "-c", "merge.autostash=false", "pull", "--ff-only", "--no-rebase", "--no-autostash")
	}
	return g.Run(ctx, dir, "-c", "rebase.autostash=false", "-c", "pull.ff=false", "pull", "--rebase", "--no-autostash")
}
