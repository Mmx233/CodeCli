package target

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Mmx233/CodeCli/internal/config"
)

func TestParse(t *testing.T) {
	cfg := config.DefaultsFor(t.TempDir())
	cfg.Default.Username = "Mmx233"
	for _, tc := range []struct{ raw, host, owner, repo, sub, scheme string }{
		{"CodeCli", "github.com", "Mmx233", "CodeCli", "", "https"},
		{"Mmx233/CodeCli", "github.com", "Mmx233", "CodeCli", "", "https"},
		{"github.com/Mmx233/CodeCli/Sub/目录", "github.com", "Mmx233", "CodeCli", filepath.Join("Sub", "目录"), "https"},
		{"Mmx233/CodeCli/Sub", "github.com", "Mmx233", "CodeCli", "Sub", "https"},
		{"https://github.com/Mmx233/CodeCli.git", "github.com", "Mmx233", "CodeCli", "", "https"},
		{"git@github.com:Mmx233/CodeCli.git", "github.com", "Mmx233", "CodeCli", "", "ssh"},
		{"ssh://git@github.com:22/Mmx233/CodeCli.git", "github.com", "Mmx233", "CodeCli", "", "ssh"},
		{"https://internal:8443/Team/Repo", "internal", "Team", "Repo", "", "https"},
	} {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := Parse(tc.raw, cfg)
			if err != nil {
				t.Fatal(err)
			}
			if got.Host != tc.host || got.Owner != tc.owner || got.Repo != tc.repo || got.Subdir != tc.sub || got.Scheme != tc.scheme {
				t.Fatalf("%+v", got)
			}
		})
	}
	for _, raw := range []string{"ext::sh dangerous", "file:///tmp/a", "ftp://x/a/b", "https://host/a", "https://host/a/b/extra", "https://github.com/a/../x", "https://github.com/a/%2e%2e", "https://github.com/a/b%2fc", "a//b", "a/b/..", "a/b/", "a/b\x00", "https://host:99999/a/b", "https://u:p@/a/b", "https://host/CON/b"} {
		if _, err := Parse(raw, cfg); err == nil {
			t.Errorf("accepted %q", raw)
		}
	}
	for _, raw := range []string{".", "..", "./folder", "../folder", filepath.Join(t.TempDir(), "missing")} {
		got, err := Parse(raw, config.Config{})
		if err != nil || !got.Local || !filepath.IsAbs(got.Path) {
			t.Fatal(raw, got, err)
		}
	}
	if runtime.GOOS == "windows" {
		for _, raw := range []string{`C:\projects\Repo`, `\\server\share\Repo`, `.\Repo`, `..\Repo`} {
			if !IsLocal(raw) {
				t.Fatal("Windows path not local", raw)
			}
		}
	}
}

func TestGitHubCaseReuseAndConflicts(t *testing.T) {
	storage := t.TempDir()
	existing := filepath.Join(storage, "GitHub.com", "Mmx233", "CodeCli")
	os.MkdirAll(existing, 0755)
	for _, raw := range []string{"github.com/MMX233/CODECLI", "https://github.com/mmx233/codecli", "git@github.com:MMX233/CodeCLI.git", "ssh://git@github.com:22/mmx233/codecli.git"} {
		tt, err := Parse(raw, config.Config{})
		if err != nil {
			t.Fatal(err)
		}
		p, err := tt.Locate(storage)
		if err != nil || p != existing {
			t.Fatal(raw, p, err)
		}
	}
	newTarget, _ := Parse("github.com/MMX233/NEW", config.Config{})
	p, err := newTarget.Locate(storage)
	if err != nil || p != filepath.Join(storage, "GitHub.com", "Mmx233", "new") {
		t.Fatal(p, err)
	}
	other, _ := Parse("https://example.com/Team/Repo", config.Config{})
	p, err = other.Locate(storage)
	if err != nil || p != filepath.Join(storage, "example.com", "Team", "Repo") {
		t.Fatal(p, err)
	}
	// APFS may fold case; use a genuinely sensitive volume for this branch.
	upper := filepath.Join(storage, "GitHub.com", "Mmx233", "CODECLI")
	if err = os.Mkdir(upper, 0755); os.IsExist(err) {
		t.Skip("case-sensitive directory conflicts require a case-sensitive filesystem")
	} else if err != nil {
		t.Fatal(err)
	}
	if _, err = newTarget.Locate(storage); err != nil {
		t.Fatal(err)
	}
	old, _ := Parse("github.com/Mmx233/CodeCli", config.Config{})
	if _, err = old.Locate(storage); err == nil {
		t.Fatal("exact match hid case conflict")
	}
}

func TestWebMappingAndRedaction(t *testing.T) {
	local, _ := Parse(t.TempDir(), config.Config{})
	if _, err := local.WebURL(); err == nil {
		t.Fatal("mapped local origin to a website")
	}
	for _, remote := range []string{"https://host/Owner/Repo", "https://host/Owner/Repo.git", "ssh://git@host/Owner/Repo", "git@host:Owner/Repo"} {
		parsed, err := Parse(remote, config.Config{})
		if err != nil || parsed.Remote != remote {
			t.Fatal("explicit remote was rewritten", remote, parsed.Remote, err)
		}
	}
	for _, tc := range []struct{ raw, want string }{
		{"https://user:secret@host:8443/Team/Repo.git?token=hidden", "https://host:8443/Team/Repo"},
		{"http://host:80/a/b", "http://host:80/a/b"},
		{"git@github.com:Owner/Repo.git", "https://github.com/Owner/Repo"},
		{"ssh://git@github.com:22/Owner/Repo.git", "https://github.com/Owner/Repo"},
	} {
		tt, err := Parse(tc.raw, config.Config{})
		if err != nil {
			t.Fatal(err)
		}
		got, err := tt.WebURL()
		if err != nil || got != tc.want {
			t.Fatal(got, err)
		}
	}
	for _, raw := range []string{"git@alias:a/b", "ssh://git@github.com:443/a/b"} {
		tt, err := Parse(raw, config.Config{})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = tt.WebURL(); err == nil {
			t.Fatal("guessed SSH URL")
		}
	}
	text := RedactText("fatal 'https://user:password@host/a/b?token=secret' failed")
	if strings.Contains(text, "password") || strings.Contains(text, "token") || strings.Contains(text, "secret") {
		t.Fatal(text)
	}
}

func TestSymlinkBoundary(t *testing.T) {
	real, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(real, "link")
	if err = os.Symlink(real, link); err != nil {
		t.Skip(err)
	}
	if err = NoSymlinks(filepath.Join(link, "new", "repo")); err == nil {
		t.Fatal("accepted symlink ancestor")
	}
	if err = NoSymlinks(filepath.Join(real, "new", "repo")); err != nil {
		t.Fatal(err)
	}
}
