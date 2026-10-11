package project

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrunePreviewAndApply(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	if err := f.git.Prune(ctx, f.storage, []string{p}, false, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("preview removed checkout")
	}
	if _, err := os.Stat(filepath.Join(f.storage, ".codecli.lock")); !os.IsNotExist(err) {
		t.Fatal("preview created lock")
	}
	if err := f.git.Prune(ctx, f.storage, []string{p}, true, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("apply did not remove checkout", err)
	}
	if _, err := os.Stat(f.seed); err != nil {
		t.Fatal("apply touched source repository")
	}
}

func TestPruneProtectsLocalData(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		change     func(*testing.T, fixture, string)
	}{
		{"dirty", "files", func(t *testing.T, f fixture, p string) {
			os.WriteFile(filepath.Join(p, "base.txt"), []byte("dirty"), 0600)
		}},
		{"untracked", "files", func(t *testing.T, f fixture, p string) {
			os.WriteFile(filepath.Join(p, "untracked"), []byte("keep"), 0600)
		}},
		{"ignored", "files", func(t *testing.T, f fixture, p string) {
			os.WriteFile(filepath.Join(p, ".git", "info", "exclude"), []byte("ignored\n"), 0600)
			os.WriteFile(filepath.Join(p, "ignored"), []byte("keep"), 0600)
		}},
		{"stash", "stash", func(t *testing.T, f fixture, p string) {
			os.WriteFile(filepath.Join(p, "base.txt"), []byte("dirty"), 0600)
			runGit(t, p, "stash")
		}},
		{"ahead", "commits", func(t *testing.T, f fixture, p string) { commitFile(t, p, "local", "local") }},
		{"local-branch-name", "branch name", func(t *testing.T, f fixture, p string) { runGit(t, p, "branch", "private-name") }},
		{"local-tag-name", "tags", func(t *testing.T, f fixture, p string) { runGit(t, p, "tag", "private-tag") }},
		{"reflog", "reflogs", func(t *testing.T, f fixture, p string) {
			commitFile(t, p, "recovery", "recovery")
			runGit(t, p, "reset", "--hard", "origin/main")
		}},
		{"detached", "detached", func(t *testing.T, f fixture, p string) { runGit(t, p, "checkout", "--detach") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			p := f.clone(t)
			tc.change(t, f, p)
			err := f.git.Prune(context.Background(), f.storage, []string{p}, true, false)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %s, got %v", tc.want, err)
			}
			if _, err = os.Stat(p); err != nil {
				t.Fatal("protected checkout removed")
			}
			if err = f.git.Prune(context.Background(), f.storage, []string{p}, false, true); err != nil {
				t.Fatal("force preview rejected data-only protection", err)
			}
			if _, err = os.Stat(p); err != nil {
				t.Fatal("force preview removed data")
			}
		})
	}
}

func TestForceNeverRelaxesBoundary(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	sub := filepath.Join(p, "Sub")
	os.Mkdir(sub, 0700)
	plain := filepath.Join(f.storage, "github.com", "owner", "ordinary")
	os.MkdirAll(plain, 0700)
	for _, bad := range []string{f.storage, filepath.Dir(p), sub, f.seed, plain} {
		if err := f.git.Prune(ctx, f.storage, []string{bad}, true, true); err == nil {
			t.Fatal("force accepted boundary", bad)
		}
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("boundary check removed checkout")
	}
	if err := f.git.Prune(ctx, f.storage, []string{p, p}, true, true); err == nil {
		t.Fatal("duplicate targets accepted")
	}
	link := filepath.Join(f.storage, "github.com", "owner", "link")
	if err := os.Symlink(p, link); err == nil {
		if err = f.git.Prune(ctx, f.storage, []string{link}, true, true); err == nil {
			t.Fatal("force followed symlink")
		}
	}
	worktree := filepath.Join(f.base, "linked")
	runGit(t, p, "worktree", "add", "-b", "other", worktree)
	if err := f.git.Prune(ctx, f.storage, []string{p}, true, true); err == nil || !strings.Contains(err.Error(), "worktrees") {
		t.Fatal("force accepted shared worktrees", err)
	}
	if _, err := os.Stat(worktree); err != nil {
		t.Fatal("external worktree changed")
	}
	if err := f.git.Prune(ctx, f.storage, []string{worktree}, true, true); err == nil {
		t.Fatal("external linked worktree accepted")
	}
}

func TestPruneExternalGitdirAndBatchValidation(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	external := filepath.Join(f.storage, "github.com", "owner", "external")
	runGit(t, f.base, "init", "--separate-git-dir", filepath.Join(f.base, "metadata"), external)
	if err := f.git.Prune(ctx, f.storage, []string{external}, true, true); err == nil {
		t.Fatal("external gitdir accepted")
	}
	if err := f.git.Prune(ctx, f.storage, []string{p, external}, true, true); err == nil {
		t.Fatal("invalid batch accepted")
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("partial batch deleted good checkout")
	}
}
