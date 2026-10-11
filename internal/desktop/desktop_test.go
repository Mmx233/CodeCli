package desktop

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"testing"

	"github.com/Mmx233/CodeCli/internal/config"
)

func has(names ...string) func(string) (string, error) {
	return func(name string) (string, error) {
		if slices.Contains(names, name) {
			return name, nil
		}
		return "", errors.New("not found")
	}
}

func TestPlatformPlans(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "空 格")
	for _, tc := range []struct {
		os, name string
		args     []string
		console  bool
	}{
		{"darwin", "", []string{"-a", "Terminal", dir}, false},
		{"linux", "gnome-terminal", []string{"--working-directory=" + dir}, false},
		{"linux", "konsole", []string{"--workdir", dir}, false},
		{"linux", "xdg-terminal-exec", []string{"--dir=" + dir}, false},
		{"windows", "powershell", nil, true},
		{"windows", "wt", []string{"-d", dir}, true},
	} {
		p, err := Terminal(tc.os, tc.name, dir, has("open", "powershell", "gnome-terminal", "konsole", "xdg-terminal-exec", "wt"))
		if err != nil || !reflect.DeepEqual(p.Args, tc.args) || p.Dir != dir || p.Console != tc.console {
			t.Fatal(tc, p, err)
		}
	}
	p, err := Terminal("linux", "", dir, has("konsole"))
	if err != nil || p.Program != "konsole" {
		t.Fatal(p, err)
	}
	if _, err = Terminal("linux", "", dir, has()); err == nil {
		t.Fatal("missing terminal accepted")
	}
	for _, goos := range []string{"linux", "darwin"} {
		if _, err = Terminal(goos, "zsh", dir, has("zsh")); err == nil {
			t.Fatal("shell mistaken for a desktop terminal")
		}
	}
	p, err = IDE("windows", "goland", dir, has("goland64.exe"))
	if err != nil || p.Program != "goland64.exe" || !reflect.DeepEqual(p.Args, []string{dir}) {
		t.Fatal(p, err)
	}
	if _, err = IDE("linux", "missing", dir, has()); err == nil {
		t.Fatal("missing IDE accepted")
	}
	for _, goos := range []string{"linux", "darwin", "windows"} {
		p, err = Web(goos, "https://github.com/a/b", has("open", "xdg-open", "rundll32"))
		if err != nil || !p.Wait {
			t.Fatal(p, err)
		}
	}
}

func TestSessionAndIDESelection(t *testing.T) {
	if err := Session("linux", func(string) string { return "" }); err == nil {
		t.Fatal("headless accepted")
	}
	if err := Session("linux", func(k string) string {
		if k == "WAYLAND_DISPLAY" {
			return "wayland-0"
		}
		return ""
	}); err != nil {
		t.Fatal(err)
	}
	if err := Session("darwin", func(k string) string {
		if k == "SSH_TTY" {
			return "/dev/ttys1"
		}
		return ""
	}); err == nil {
		t.Fatal("headless SSH accepted")
	}
	dir := t.TempDir()
	cfg := config.DefaultsFor(dir)
	cfg.Default.Idea = "default"
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0600)
	for _, tc := range []struct{ explicit, want string }{{"chosen", "chosen"}, {"", "goland"}} {
		got, err := SelectIDE(cfg, tc.explicit, dir)
		if err != nil || got != tc.want {
			t.Fatal(got, err)
		}
	}
	os.Remove(filepath.Join(dir, "go.mod"))
	got, err := SelectIDE(cfg, "", dir)
	if err != nil || got != "default" {
		t.Fatal(got, err)
	}
	cfg.Default.Idea = ""
	if _, err = SelectIDE(cfg, "", dir); err == nil {
		t.Fatal("empty default accepted")
	}
	cfg.Rules = []config.Rule{{Idea: "wrong", File: []string{"../outside"}}}
	if _, err = SelectIDE(cfg, "", dir); err == nil {
		t.Fatal("escaping rule accepted")
	}
}

func TestLauncherArgvCwdAndFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX launcher fixture; Windows native launch still needs an environment")
	}
	dir := filepath.Join(t.TempDir(), "空 格")
	os.Mkdir(dir, 0700)
	launcher := filepath.Join(dir, "launcher")
	if err := os.WriteFile(launcher, []byte("#!/bin/sh\npwd -P\nprintf '%s\\n' \"$@\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	p := Plan{Program: launcher, Args: []string{"argument with space", "中文"}, Dir: dir, Wait: true}
	if err := p.Launch(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	real, _ := filepath.EvalSymlinks(dir)
	if output.String() != real+"\nargument with space\n中文\n" {
		t.Fatal(output.String())
	}
	os.WriteFile(launcher, []byte("#!/bin/sh\nexit 23\n"), 0700)
	if err := p.Launch(context.Background(), &output); err == nil {
		t.Fatal("launcher failure swallowed")
	}
	p.Program = filepath.Join(dir, "missing")
	p.Wait = false
	if err := p.Launch(context.Background(), &output); err == nil {
		t.Fatal("start failure swallowed")
	}
}
