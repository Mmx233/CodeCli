package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Mmx233/CodeCli/internal/config"
)

func execute(store func() (config.Store, error), args ...string) (string, string, error) {
	c := newRoot("test-version", store)
	var out, stderr bytes.Buffer
	c.SetArgs(args)
	c.SetOut(&out)
	c.SetErr(&stderr)
	c.SetIn(strings.NewReader(""))
	err := c.Execute()
	return out.String(), stderr.String(), err
}

func TestHelpVersionCompletionDoNotLoadConfig(t *testing.T) {
	store := func() (config.Store, error) {
		t.Fatal("config loaded")
		return config.Store{}, errors.New("should not load")
	}
	for _, args := range [][]string{{"--help"}, {"-v"}, {"--version"}, {"open", "--help"}, {"clone", "--help"}, {"config", "--help"}, {"prune", "--help"}, {"completion", "bash"}, {"completion", "zsh"}, {"completion", "fish"}, {"completion", "powershell"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			out, _, err := execute(store, args...)
			if err != nil || out == "" {
				t.Fatal(out, err)
			}
		})
	}
}

func TestArgumentErrorsHaveNoSideEffects(t *testing.T) {
	store := func() (config.Store, error) { t.Fatal("config loaded for invalid args"); return config.Store{}, nil }
	for _, args := range [][]string{
		{"cmd", "Repo"}, {"browser", "Repo"}, {"clear", "Repo"}, {"project", "Repo"}, {"path"}, {"version"},
		{"open", "--web", "--terminal"}, {"open", "--idea", "goland"}, {"--idea", "goland"},
		{"clone"}, {"clone", "x", "--fetch", "--pull"}, {"clone", "x", "--pull", "--rebase"},
		{"clone", "x", "--web"}, {"--web", "clone", "x"}, {"pwd", "x", "--terminal"}, {"prune"},
		{"config", "set", "default.username=old"}, {"config", "set", "rules"}, {"config", "get", "x", "y"},
	} {
		out, _, err := execute(store, args...)
		if err == nil || out != "" {
			t.Fatalf("%v => %q, %v", args, out, err)
		}
	}
	_, _, err := execute(store, "cmd", "Repo")
	if !strings.Contains(err.Error(), "--terminal") {
		t.Fatal(err)
	}
}

func TestPwdFilesystemOnly(t *testing.T) {
	noConfig := func() (config.Store, error) { t.Fatal("local pwd loaded config"); return config.Store{}, nil }
	cwd, _ := os.Getwd()
	out, stderr, err := execute(noConfig, "pwd")
	if err != nil || out != cwd+"\n" || stderr != "" {
		t.Fatal(out, stderr, err)
	}
	dir := filepath.Join(t.TempDir(), "空 格")
	os.Mkdir(dir, 0700)
	t.Setenv("PATH", "")
	out, stderr, err = execute(noConfig, "pwd", dir)
	if err != nil || out != dir+"\n" || stderr != "" {
		t.Fatal(out, stderr, err)
	}
	out, _, err = execute(noConfig, "pwd", filepath.Join(dir, "missing"))
	if err == nil || out != "" {
		t.Fatal(out, err)
	}
	configPath := filepath.Join(t.TempDir(), "config")
	store := func() (config.Store, error) { return config.Store{Path: configPath}, nil }
	out, _, err = execute(store, "pwd", "Owner/Repo")
	if err == nil || out != "" {
		t.Fatal(out, err)
	}
	if _, err = os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatal("pwd created config")
	}
}

func TestConfigPercentAndCanonicalSyntax(t *testing.T) {
	name := filepath.Join(t.TempDir(), "config")
	store := func() (config.Store, error) { return config.Store{Path: name}, nil }
	if out, _, err := execute(store, "config", "get"); err != nil || !strings.Contains(out, "gitSite: github.com") {
		t.Fatal(out, err)
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Fatal("get wrote config")
	}
	if out, _, err := execute(store, "config", "set", "default.username", "name=100%"); err != nil || out != "" {
		t.Fatal(out, err)
	}
	out, _, err := execute(store, "config", "get", "default.username")
	if err != nil || out != "name=100%\n" {
		t.Fatal(out, err)
	}
	out, _, err = execute(store, "config", "get")
	if err != nil || strings.Contains(out, "%!") || !strings.Contains(out, "name=100%") {
		t.Fatal(out, err)
	}
	if _, _, err = execute(store, "config", "unset", "default.username"); err != nil {
		t.Fatal(err)
	}
	if out, _, err = execute(store, "config", "get", "default.username"); err != nil || out != "\n" {
		t.Fatal(out, err)
	}
}

func TestCommandClonePruneAndLauncherPreflight(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	gitcfg := filepath.Join(base, "gitconfig")
	os.WriteFile(gitcfg, []byte("[user]\n name = Test\n email = test@example.invalid\n[init]\n defaultBranch = main\n[commit]\n gpgsign = false\n"), 0600)
	t.Setenv("GIT_CONFIG_GLOBAL", gitcfg)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
	git := func(dir string, args ...string) string {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dir
		b, e := c.CombinedOutput()
		if e != nil {
			t.Fatalf("git %v: %v %s", args, e, b)
		}
		return strings.TrimSpace(string(b))
	}
	source := filepath.Join(base, "source")
	git(base, "init", source)
	os.WriteFile(filepath.Join(source, "base.txt"), []byte("data"), 0600)
	git(source, "add", ".")
	git(source, "commit", "-m", "initial")
	// Exercise the CLI's HTTPS parser and let Git's own insteadOf choose a local remote.
	git(base, "config", "--file", gitcfg, "url."+source+".insteadOf", "https://github.com/Owner/Repo.git")
	storage := filepath.Join(base, "storage")
	cfgPath := filepath.Join(base, "config")
	store := func() (config.Store, error) { return config.Store{Path: cfgPath}, nil }
	if err = (config.Store{Path: cfgPath}).Set("storage.projectDir", &storage); err != nil {
		t.Fatal(err)
	}
	out, stderr, err := execute(store, "clone", "Owner/Repo", "--pull")
	p := filepath.Join(storage, "github.com", "owner", "repo")
	if err != nil || out != p+"\n" || !strings.Contains(stderr, "Cloning into") {
		t.Fatal(out, stderr, err)
	}
	git(p, "remote", "set-url", "origin", "https://github.com/Fork/Different.git")
	out, _, err = execute(store, "pwd", "OWNER/REPO")
	if err != nil || out != p+"\n" {
		t.Fatal("pwd compared origin", out, err)
	}
	out, _, err = execute(store, "clone", "owner/repo")
	if err != nil || out != p+"\n" {
		t.Fatal("reuse compared origin", out, err)
	}
	out, _, err = execute(store, "prune", "Owner/Repo")
	if err != nil || !strings.Contains(out, "Would remove") {
		t.Fatal(out, err)
	}
	if _, err = os.Stat(p); err != nil {
		t.Fatal("preview deleted")
	}
	out, _, err = execute(store, "clone", "Owner/Missing")
	if err == nil || out != "" {
		t.Fatal(out, err)
	}
	if runtime.GOOS == "linux" {
		t.Setenv("DISPLAY", ":test")
	}
	t.Setenv("SSH_CONNECTION", "")
	t.Setenv("SSH_TTY", "")
	for _, args := range [][]string{{"open", "Owner/NeverClone", "--ide", "definitely-missing-codecli-ide"}, {"Owner/NeverClone", "--ide", "definitely-missing-codecli-ide"}} {
		out, _, err = execute(store, args...)
		if err == nil || out != "" {
			t.Fatal(args, out, err)
		}
		if _, err = os.Stat(filepath.Join(storage, "github.com", "owner", "neverclone")); !os.IsNotExist(err) {
			t.Fatal("preflight failure cloned")
		}
	}
	out, _, err = execute(store, "prune", "Owner/Repo", "--apply")
	if err != nil || !strings.Contains(out, "Removed") {
		t.Fatal(out, err)
	}
	if _, err = os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("apply did not remove")
	}
	out, _, err = execute(store, "open", "Owner/Repo")
	if err == nil || out != "" {
		t.Fatal("automatic IDE selection should fail for generic repo with no default", out, err)
	}
	if _, err = os.Stat(filepath.Join(p, "base.txt")); err != nil {
		t.Fatal("IDE failure removed successful clone", err)
	}
}

func TestWebFailureIsReadOnlyAndRedactsCredentials(t *testing.T) {
	base := t.TempDir()
	name := filepath.Join(base, "config")
	store := func() (config.Store, error) { return config.Store{Path: name}, nil }
	t.Setenv("DISPLAY", "")
	t.Setenv("WAYLAND_DISPLAY", "")
	t.Setenv("SSH_CONNECTION", "test-session")
	t.Setenv("PATH", "")
	out, stderr, err := execute(store, "open", "https://user:password@host:8443/Owner/Repo.git?token=secret", "--web")
	if err == nil || out != "" || !strings.Contains(stderr, "https://host:8443/Owner/Repo") {
		t.Fatal(out, stderr, err)
	}
	if strings.Contains(stderr, "password") || strings.Contains(stderr, "secret") || strings.Contains(stderr, "token=") {
		t.Fatal(stderr)
	}
	for _, p := range []string{name, filepath.Join(base, "project")} {
		if _, err = os.Stat(p); !os.IsNotExist(err) {
			t.Fatal("web wrote config or cloned", p, err)
		}
	}
}

func TestOpenLaunchesWithCorrectArguments(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable fixture; platform argv plans have separate tests")
	}
	base := t.TempDir()
	bin := filepath.Join(base, "bin")
	os.Mkdir(bin, 0700)
	record := filepath.Join(base, "record")
	t.Setenv("CODECLI_TEST_LAUNCH_RECORD", record)
	t.Setenv("DISPLAY", ":test")
	t.Setenv("SSH_CONNECTION", "")
	t.Setenv("SSH_TTY", "")
	launcher := filepath.Join(bin, "launcher")
	script := []byte("#!/bin/sh\npwd -P > \"$CODECLI_TEST_LAUNCH_RECORD\"\nprintf '%s\\n' \"$@\" >> \"$CODECLI_TEST_LAUNCH_RECORD\"\n")
	os.WriteFile(launcher, script, 0700)
	store := func() (config.Store, error) { return config.Store{Path: filepath.Join(base, "config")}, nil }
	dir := filepath.Join(base, "空 格")
	os.Mkdir(dir, 0700)
	out, _, err := execute(store, "open", dir, "--ide", launcher)
	if err != nil || out != "" {
		t.Fatal(out, err)
	}
	deadline := time.After(3 * time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	real, _ := filepath.EvalSymlinks(dir)
	for {
		data, _ := os.ReadFile(record)
		if string(data) == real+"\n"+dir+"\n" {
			break
		}
		select {
		case <-deadline:
			t.Fatal("detached launcher did not receive argv/cwd", string(data))
		case <-tick.C:
		}
	}
	os.Remove(record)
	webName := "xdg-open"
	if runtime.GOOS == "darwin" {
		webName = "open"
	}
	os.WriteFile(filepath.Join(bin, webName), script, 0700)
	t.Setenv("PATH", bin)
	out, _, err = execute(store, "https://user:password@github.com/Owner/Repo.git?token=secret", "--web")
	if err != nil || out != "" {
		t.Fatal(out, err)
	}
	data, _ := os.ReadFile(record)
	if !strings.HasSuffix(string(data), "https://github.com/Owner/Repo\n") || strings.Contains(string(data), "password") || strings.Contains(string(data), "secret") {
		t.Fatal(string(data))
	}
	if _, err = os.Stat(filepath.Join(base, "project")); !os.IsNotExist(err) {
		t.Fatal("web cloned")
	}
}
