package project

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Mmx233/CodeCli/internal/target"
)

func isolatedGit(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	name := filepath.Join(dir, "gitconfig")
	os.WriteFile(name, []byte("[user]\n name = Test\n email = test@example.invalid\n[init]\n defaultBranch = main\n[commit]\n gpgsign = false\n[tag]\n gpgsign = false\n"), 0600)
	t.Setenv("GIT_CONFIG_GLOBAL", name)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
	for _, key := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_AUTHOR_DATE", "GIT_COMMITTER_DATE"} {
		t.Setenv(key, "")
	}
	return dir
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GIT_DIR=") && !strings.HasPrefix(e, "GIT_WORK_TREE=") && !strings.HasPrefix(e, "GIT_INDEX_FILE=") && !strings.HasPrefix(e, "GIT_COMMON_DIR=") && !strings.HasPrefix(e, "GIT_OBJECT_DIRECTORY=") && !strings.HasPrefix(e, "GIT_ALTERNATE_OBJECT_DIRECTORIES=") {
			env = append(env, e)
		}
	}
	c.Env = env
	b, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, b)
	}
	return strings.TrimSpace(string(b))
}

type fixture struct {
	base, remote, seed, storage string
	remoteTarget                target.Target
	git                         Git
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	base := isolatedGit(t)
	f := fixture{base: base, remote: filepath.Join(base, "remote.git"), seed: filepath.Join(base, "seed"), storage: filepath.Join(base, "storage"), git: Git{}}
	runGit(t, base, "init", "--bare", f.remote)
	runGit(t, base, "init", f.seed)
	os.WriteFile(filepath.Join(f.seed, "base.txt"), []byte("base\n"), 0600)
	runGit(t, f.seed, "add", ".")
	runGit(t, f.seed, "commit", "-m", "initial")
	runGit(t, f.seed, "remote", "add", "origin", f.remote)
	runGit(t, f.seed, "push", "-u", "origin", "main")
	f.remoteTarget = target.Target{Host: "github.com", Owner: "Owner", Repo: "Repo", Scheme: "https", Remote: f.remote}
	return f
}

func (f fixture) clone(t *testing.T) string {
	t.Helper()
	p, err := f.git.Clone(context.Background(), f.remoteTarget, f.storage, "")
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func commitFile(t *testing.T, dir, name, data string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", name)
	runGit(t, dir, "commit", "-m", name)
}

func TestCloneReuseOutputAndFailureOwnership(t *testing.T) {
	f := newFixture(t)
	var stderr bytes.Buffer
	f.git.Err = &stderr
	trace := filepath.Join(f.base, "trace")
	t.Setenv("GIT_TRACE", trace)
	p, err := f.git.Clone(context.Background(), f.remoteTarget, f.storage, "rebase")
	if err != nil {
		t.Fatal(err)
	}
	if p != filepath.Join(f.storage, "github.com", "owner", "repo") {
		t.Fatal(p)
	}
	data, _ := os.ReadFile(trace)
	if strings.Contains(string(data), " pull ") || strings.Contains(string(data), " fetch ") {
		t.Fatal("fresh clone was updated again", string(data))
	}
	if !strings.Contains(stderr.String(), "Cloning into") {
		t.Fatal("missing Git stderr", stderr.String())
	}
	before := runGit(t, p, "rev-parse", "HEAD")
	commitFile(t, f.seed, "remote.txt", "remote")
	runGit(t, f.seed, "push")
	if got, e := f.git.Clone(context.Background(), f.remoteTarget, f.storage, ""); e != nil || got != p {
		t.Fatal(got, e)
	}
	if after := runGit(t, p, "rev-parse", "HEAD"); after != before {
		t.Fatal("no-flag reuse changed HEAD")
	}
	missing := f.remoteTarget
	missing.Repo = "Missing"
	missing.Remote = filepath.Join(f.base, "not-a-remote")
	if got, e := f.git.Clone(context.Background(), missing, f.storage, ""); e == nil || got != "" {
		t.Fatal(got, e)
	}
	if _, e := os.Stat(p); e != nil {
		t.Fatal("failed clone removed existing repo", e)
	}
	stages, _ := filepath.Glob(filepath.Join(filepath.Dir(p), ".codecli-clone-*"))
	if len(stages) != 0 {
		t.Fatal(stages)
	}
	plain := filepath.Join(f.storage, "github.com", "owner", "plain")
	os.MkdirAll(plain, 0755)
	os.WriteFile(filepath.Join(plain, "keep"), []byte("keep"), 0600)
	other := f.remoteTarget
	other.Repo = "plain"
	if _, e := f.git.Clone(context.Background(), other, f.storage, ""); e == nil {
		t.Fatal("ordinary folder accepted as clone")
	}
	if _, e := os.Stat(filepath.Join(plain, "keep")); e != nil {
		t.Fatal("ordinary folder lost")
	}
}

func TestFetchPullAndRebase(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	initial := runGit(t, p, "rev-parse", "HEAD")
	commitFile(t, f.seed, "remote.txt", "remote")
	runGit(t, f.seed, "push")
	if err := f.git.Update(ctx, p, "fetch"); err != nil {
		t.Fatal(err)
	}
	if head := runGit(t, p, "rev-parse", "HEAD"); head != initial {
		t.Fatal("fetch moved HEAD")
	}
	if head := runGit(t, p, "rev-parse", "origin/main"); head == initial {
		t.Fatal("fetch did not update refs")
	}
	if err := f.git.Update(ctx, p, "pull"); err != nil {
		t.Fatal(err)
	}
	if runGit(t, p, "rev-parse", "HEAD") != runGit(t, f.seed, "rev-parse", "HEAD") {
		t.Fatal("pull did not fast-forward")
	}
	commitFile(t, p, "local.txt", "local")
	local := runGit(t, p, "rev-parse", "HEAD")
	commitFile(t, f.seed, "next.txt", "next")
	runGit(t, f.seed, "push")
	if err := f.git.Update(ctx, p, "pull"); err == nil {
		t.Fatal("ff-only merged divergent branches")
	}
	if runGit(t, p, "rev-parse", "HEAD") != local {
		t.Fatal("failed pull changed local HEAD")
	}
	runGit(t, p, "config", "pull.ff", "only")
	runGit(t, p, "config", "rebase.autostash", "true")
	if err := f.git.Update(ctx, p, "rebase"); err != nil {
		t.Fatal(err)
	}
	if runGit(t, p, "rev-parse", "HEAD") == local {
		t.Fatal("rebase did not rewrite local commit")
	}
	if _, err := os.Stat(filepath.Join(p, "local.txt")); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateRejectsUnsafeStatesAndPreservesConflict(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	os.WriteFile(filepath.Join(p, "untracked"), []byte("keep"), 0600)
	runGit(t, p, "config", "rebase.autostash", "true")
	for _, mode := range []string{"pull", "rebase"} {
		if err := f.git.Update(ctx, p, mode); err == nil {
			t.Fatal("accepted dirty", mode)
		}
	}
	if runGit(t, p, "stash", "list") != "" {
		t.Fatal("autostashed")
	}
	os.Remove(filepath.Join(p, "untracked"))
	runGit(t, p, "checkout", "--detach")
	if err := f.git.Update(ctx, p, "pull"); err == nil || !strings.Contains(err.Error(), "detached") {
		t.Fatal(err)
	}
	runGit(t, p, "checkout", "main")
	runGit(t, p, "branch", "--unset-upstream")
	if err := f.git.Update(ctx, p, "rebase"); err == nil || !strings.Contains(err.Error(), "upstream") {
		t.Fatal(err)
	}
	runGit(t, p, "branch", "--set-upstream-to=origin/main")
	os.Mkdir(filepath.Join(p, ".git", "rebase-apply"), 0700)
	if err := f.git.Update(ctx, p, "pull"); err == nil || !strings.Contains(err.Error(), "in progress") {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(p, ".git", "rebase-apply"))
	commitFile(t, p, "base.txt", "local\n")
	commitFile(t, f.seed, "base.txt", "remote\n")
	runGit(t, f.seed, "push")
	if err := f.git.Update(ctx, p, "rebase"); err == nil {
		t.Fatal("conflicting rebase succeeded")
	}
	if !strings.Contains(runGit(t, p, "status", "--porcelain"), "UU") {
		t.Fatal("conflict not preserved")
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatal("conflict removed checkout")
	}
}

func TestConcurrentCloneAndNoReplace(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	var wg sync.WaitGroup
	start := make(chan struct{})
	out := make(chan error, 2)
	for range 2 {
		wg.Go(func() { <-start; _, err := f.git.Clone(ctx, f.remoteTarget, f.storage, ""); out <- err })
	}
	close(start)
	wg.Wait()
	close(out)
	success := 0
	for err := range out {
		if err == nil {
			success++
		} else if !strings.Contains(err.Error(), "locked") {
			t.Fatal(err)
		}
	}
	if success == 0 {
		t.Fatal("neither clone succeeded")
	}
	p, _ := f.remoteTarget.Locate(f.storage)
	if err := f.git.CheckRoot(ctx, p); err != nil {
		t.Fatal("winner checkout lost", err)
	}
	source := filepath.Join(f.base, "stage")
	destination := filepath.Join(f.base, "destination")
	os.Mkdir(source, 0700)
	os.Mkdir(destination, 0700)
	os.WriteFile(filepath.Join(source, "owned"), []byte("own"), 0600)
	if err := publish(source, destination); err == nil {
		t.Fatal("replaced an existing empty destination")
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatal("source removed on publication failure")
	}
	if _, err := os.Stat(destination); err != nil {
		t.Fatal("destination removed")
	}
}

func TestManagedUpdateHonorsPruneLock(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	unlock, err := lock(f.storage)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if _, err = f.git.Clone(context.Background(), f.remoteTarget, f.storage, "fetch"); err == nil || !strings.Contains(err.Error(), "locked") {
		t.Fatal("managed sync bypassed prune lock", err)
	}
	if _, err = os.Stat(p); err != nil {
		t.Fatal("busy checkout lost")
	}
}

func TestLocalCloneAndSubdirectoryOpen(t *testing.T) {
	f := newFixture(t)
	p := f.clone(t)
	ctx := context.Background()
	tt := f.remoteTarget
	tt.Subdir = "missing"
	if _, err := f.git.Clone(ctx, tt, f.storage, ""); err == nil {
		t.Fatal("cloned subdir")
	}
	if _, err := f.git.Ensure(ctx, tt, f.storage); err == nil {
		t.Fatal("opened missing subdir")
	}
	os.Mkdir(filepath.Join(p, "Sub"), 0700)
	tt.Subdir = "Sub"
	if got, err := f.git.Ensure(ctx, tt, f.storage); err != nil || got != filepath.Join(p, "Sub") {
		t.Fatal(got, err)
	}
	local := target.Target{Local: true, Path: p}
	if got, err := f.git.Clone(ctx, local, "", "fetch"); err != nil || got != p {
		t.Fatal(got, err)
	}
	local.Path = filepath.Join(f.base, "absent")
	if _, err := f.git.Clone(ctx, local, "", ""); err == nil {
		t.Fatal("missing local cloned")
	}
	t.Setenv("GIT_DIR", filepath.Join(f.seed, ".git"))
	t.Setenv("GIT_WORK_TREE", f.seed)
	if err := f.git.CheckRoot(ctx, p); err != nil {
		t.Fatal("Git env redirected root check", err)
	}
}
