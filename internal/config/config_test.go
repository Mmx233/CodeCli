package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReadOnlyDefaultsAndPreservingEdits(t *testing.T) {
	s := Store{Path: filepath.Join(t.TempDir(), ".CodeCli.yaml")}
	cfg, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Default.GitSite != "github.com" || len(cfg.Rules) == 0 {
		t.Fatal(cfg)
	}
	if _, err = os.Stat(s.Path); !os.IsNotExist(err) {
		t.Fatal("load wrote config", err)
	}
	input := `# my settings
default:
  gitSite: github.com
  username: Old # keep this comment
  idea: goland
  unknown: 42
storage:
  projectDir: /some/path
rules: []
extra:
  value: 100%
`
	if err = os.WriteFile(s.Path, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	value := "New=100%"
	if err = s.Set("default.username", &value); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	for _, text := range []string{"# my settings", "# keep this comment", "unknown: 42", "value: 100%", "idea: goland", "projectDir: /some/path"} {
		if !strings.Contains(string(data), text) {
			t.Fatalf("lost %q: %s", text, data)
		}
	}
	cfg, err = s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Default.Username != value {
		t.Fatal(cfg)
	}
	if err = s.Set("default.username", nil); err != nil {
		t.Fatal(err)
	}
	cfg, err = s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Default.Username != "" || cfg.Default.Idea != "goland" {
		t.Fatal(cfg)
	}
}

func TestInvalidConfigNeverChangesOriginal(t *testing.T) {
	for _, data := range []string{
		"default: {username: 42}\n", "rules: [{idea: test, file: nope}]\n",
		"default: {username: a, username: b}\n", "default: {username: a}\n---\nextra: x\n",
		"default: {username: &owner original}\nextra: *owner\n",
		"default: &section {username: original}\nextra: *section\n",
		"default: {username: [a]}\n", "default: &defaults {username: a}\nstorage: {<<: *defaults}\n",
	} {
		t.Run(data, func(t *testing.T) {
			s := Store{Path: filepath.Join(t.TempDir(), "config")}
			os.WriteFile(s.Path, []byte(data), 0600)
			v := "new"
			if err := s.Set("default.username", &v); err == nil {
				t.Fatal("accepted invalid YAML")
			}
			got, _ := os.ReadFile(s.Path)
			if string(got) != data {
				t.Fatal("original changed")
			}
		})
	}
	s := Store{Path: filepath.Join(t.TempDir(), "config")}
	v := "x"
	for _, key := range []string{"rules", "rules.idea", "default", "Default.Username", "unknown"} {
		if err := s.Set(key, &v); err == nil {
			t.Fatal("accepted", key)
		}
	}
	if _, err := os.Stat(s.Path); !os.IsNotExist(err) {
		t.Fatal("invalid key wrote config")
	}
}

func TestWriteFailurePreservesExistingConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX directory permissions; Windows sharing test covers replacement failure")
	}
	dir := t.TempDir()
	s := Store{Path: filepath.Join(dir, "config")}
	original := []byte("default: {username: old}\n")
	if err := os.WriteFile(s.Path, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0500); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0700)
	v := "new"
	if err := s.Set("default.username", &v); err == nil {
		t.Fatal("write to unwritable directory succeeded")
	}
	got, err := os.ReadFile(s.Path)
	if err != nil || string(got) != string(original) {
		t.Fatal("original changed after write failure", err)
	}
}

func TestAtomicWriteFailureAndReplacement(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "config")
	if err := atomicWrite(name, []byte("old")); err != nil {
		t.Fatal(err)
	}
	if err := atomicWrite(name, []byte("new")); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(name)
	if string(b) != "new" {
		t.Fatal(string(b))
	}
	blocked := filepath.Join(dir, "directory")
	os.Mkdir(blocked, 0700)
	os.WriteFile(filepath.Join(blocked, "keep"), []byte("keep"), 0600)
	if err := atomicWrite(blocked, []byte("bad")); err == nil {
		t.Fatal("replacement should fail")
	}
	b, err := os.ReadFile(filepath.Join(blocked, "keep"))
	if err != nil || string(b) != "keep" {
		t.Fatal("replacement removed old destination", err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, ".codecli-config-*"))
	if len(files) != 0 {
		t.Fatal("leaked staging files", files)
	}
	link := filepath.Join(dir, "link")
	if err = os.Symlink(name, link); err != nil {
		t.Skip("symlinks unavailable", err)
	}
	v := "x"
	if err = (Store{Path: link}).Set("default.username", &v); err == nil {
		t.Fatal("edited symlink")
	}
	b, _ = os.ReadFile(name)
	if string(b) != "new" {
		t.Fatal("symlink target changed")
	}
}
