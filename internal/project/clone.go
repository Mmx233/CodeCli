package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mmx233/CodeCli/internal/target"
)

// Clone and prune serialize only managed writes. A stale lock is never stolen.
func lock(storage string) (func(), error) {
	if err := target.NoSymlinks(storage); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(storage, 0755); err != nil {
		return nil, err
	}
	p := filepath.Join(storage, ".codecli.lock")
	if err := os.Mkdir(p, 0700); err != nil {
		return nil, fmt.Errorf("managed storage is locked (%s); another clone/prune may be running; remove a stale lock only after checking: %w", p, err)
	}
	return func() { _ = os.Remove(p) }, nil
}

func (g Git) Clone(ctx context.Context, t target.Target, storage, mode string) (string, error) {
	if t.Subdir != "" {
		return "", errors.New("clone expects a repository root; use open for subdirectories")
	}
	if t.Local {
		if _, err := target.ExistingDir(t.Path); err != nil {
			return "", err
		}
		if err := g.Update(ctx, t.Path, mode); err != nil {
			return "", err
		}
		return t.Path, nil
	}
	p, err := t.Locate(storage)
	if err != nil {
		return "", err
	}
	if err = target.NoSymlinks(p); err != nil {
		return "", err
	}
	unlock, err := lock(storage)
	if err != nil {
		return "", err
	}
	defer unlock()
	// Re-resolve after locking: a concurrent clone or prune may have finished.
	p, err = t.Locate(storage)
	if err != nil {
		return "", err
	}
	if err = target.NoSymlinks(p); err != nil {
		return "", err
	}
	if info, e := os.Stat(p); e == nil {
		if !info.IsDir() {
			return "", errors.New("checkout path is not a directory")
		}
		if err = g.Update(ctx, p, mode); err != nil {
			return "", err
		}
		return p, nil
	} else if !errors.Is(e, os.ErrNotExist) {
		return "", e
	}
	parent := filepath.Dir(p)
	if err = os.MkdirAll(parent, 0755); err != nil {
		return "", err
	}
	stage, err := os.MkdirTemp(parent, ".codecli-clone-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(stage) // Only the uniquely owned staging directory.
	if err = g.Run(ctx, "", "clone", "--", t.Remote, stage); err != nil {
		return "", err
	}
	if err = target.NoSymlinks(p); err != nil {
		return "", err
	}
	if err = publish(stage, p); err != nil {
		return "", fmt.Errorf("publish checkout failed; existing destination preserved: %w", err)
	}
	return p, nil
}

func (g Git) Ensure(ctx context.Context, t target.Target, storage string) (string, error) {
	if t.Local {
		return target.ExistingDir(t.Path)
	}
	root, err := t.Locate(storage)
	if err != nil {
		return "", err
	}
	if _, e := target.ExistingDir(root); errors.Is(e, os.ErrNotExist) {
		whole := t
		whole.Subdir = ""
		root, err = g.Clone(ctx, whole, storage, "")
		if err != nil {
			return "", err
		}
	} else if e != nil {
		return "", e
	}
	return target.ExistingDir(filepath.Join(root, t.Subdir))
}
