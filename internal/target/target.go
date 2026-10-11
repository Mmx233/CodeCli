package target

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/Mmx233/CodeCli/internal/config"
)

type Target struct {
	Local             bool
	Path              string
	Host, Owner, Repo string
	Subdir            string
	Remote            string
	Scheme            string
	Port              string
}

func IsLocal(raw string) bool {
	return raw == "" || raw == "." || raw == ".." || filepath.IsAbs(raw) ||
		strings.HasPrefix(raw, "./") || strings.HasPrefix(raw, "../") ||
		(filepath.Separator == '\\' && (strings.HasPrefix(raw, `.\`) || strings.HasPrefix(raw, `..\`)))
}

func Parse(raw string, cfg config.Config) (Target, error) {
	if IsLocal(raw) {
		if raw == "" {
			raw = "."
		}
		p, err := filepath.Abs(raw)
		if strings.ContainsAny(p, "\x00\r\n") {
			return Target{}, errors.New("path contains control characters")
		}
		return Target{Local: true, Path: p}, err
	}
	t := Target{Scheme: "https"}
	var parts []string
	var u *url.URL
	var err error
	if strings.Contains(raw, "://") {
		u, err = url.Parse(raw)
		if err != nil {
			return t, errors.New("invalid remote URL")
		}
		if u.Scheme != "https" && u.Scheme != "http" && u.Scheme != "ssh" {
			return t, errors.New("supported remote protocols are http, https and ssh")
		}
		t.Scheme = u.Scheme
		t.Host = u.Host
		t.Port = u.Port()
		if u.Opaque != "" || u.Host == "" {
			return t, errors.New("remote URL needs a host")
		}
		for p := range strings.SplitSeq(strings.TrimPrefix(u.EscapedPath(), "/"), "/") {
			v, e := url.PathUnescape(p)
			if e != nil {
				return t, errors.New("invalid remote path escaping")
			}
			parts = append(parts, v)
		}
	} else if at := strings.Index(raw, "@"); at >= 0 {
		colon := strings.Index(raw[at+1:], ":")
		if colon < 0 {
			return t, errors.New("SCP remote must use user@host:owner/repo")
		}
		colon += at + 1
		if raw[:at] == "" || strings.ContainsAny(raw[:at], "/\\\r\n ") {
			return t, errors.New("invalid SSH user")
		}
		t.Scheme = "ssh"
		t.Host = raw[at+1 : colon]
		parts = strings.Split(raw[colon+1:], "/")
		t.Remote = raw[:colon+1]
	} else {
		if strings.Contains(raw, "::") || strings.Contains(raw, "\\") {
			return t, errors.New("unsupported remote syntax; use an explicit URL or local path")
		}
		parts = strings.Split(raw, "/")
		switch {
		case len(parts) >= 3 && strings.Contains(parts[0], "."):
			t.Host = parts[0]
			parts = parts[1:]
		case len(parts) == 1:
			t.Host = cfg.Default.GitSite
			if cfg.Default.Username == "" {
				return t, errors.New("repository shorthand needs default.username; use owner/repo")
			}
			parts = append([]string{cfg.Default.Username}, parts...)
		default:
			t.Host = cfg.Default.GitSite
		}
	}
	host, port, err := hostParts(t.Host)
	if err != nil {
		return t, err
	}
	t.Host = host
	t.Port = port
	if len(parts) < 2 {
		return t, errors.New("remote needs owner/repo")
	}
	remoteRepo := parts[1]
	parts[1] = strings.TrimSuffix(parts[1], ".git")
	for _, part := range parts {
		if err = ValidSegment(part); err != nil {
			return t, err
		}
	}
	if len(parts) > 2 && t.Host != "github.com" {
		return t, errors.New("nested namespaces or remote subdirectories on this host are unsupported; use a local path")
	}
	t.Owner, t.Repo = parts[0], parts[1]
	if len(parts) > 2 {
		t.Subdir = filepath.Join(parts[2:]...)
	}
	if u != nil {
		u.Path = "/" + t.Owner + "/" + remoteRepo
		u.RawPath = ""
		u.Fragment = ""
		t.Remote = u.String()
	} else if t.Scheme == "ssh" {
		t.Remote += t.Owner + "/" + remoteRepo
	} else {
		t.Remote = (&url.URL{Scheme: "https", Host: t.HostPort(), Path: "/" + t.Owner + "/" + t.Repo + ".git"}).String()
	}
	return t, nil
}

func hostParts(raw string) (string, string, error) {
	u, err := url.Parse("https://" + raw)
	if err != nil || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return "", "", errors.New("invalid Git host")
	}
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`).MatchString(host) || strings.Contains(host, "..") {
		return "", "", errors.New("unsupported Git host (IPv6 literals are not supported)")
	}
	if port != "" {
		n, e := strconv.Atoi(port)
		if e != nil || n < 1 || n > 65535 {
			return "", "", errors.New("invalid Git port")
		}
	}
	if err = ValidSegment(host); err != nil {
		return "", "", err
	}
	return host, port, nil
}

func ValidSegment(s string) error {
	if s == "" || s == "." || s == ".." || strings.ContainsAny(s, "/\\:*?\"<>|%") || strings.TrimSpace(s) != s || strings.HasSuffix(s, ".") {
		return errors.New("invalid repository identity or subdirectory segment")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return errors.New("remote path contains control characters")
		}
	}
	base := strings.ToUpper(strings.SplitN(s, ".", 2)[0])
	if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || regexp.MustCompile(`^(COM|LPT)[1-9]$`).MatchString(base) {
		return errors.New("remote path uses a reserved filesystem name")
	}
	return nil
}

func (t Target) HostPort() string {
	if t.Port != "" {
		return t.Host + ":" + t.Port
	}
	return t.Host
}
func (t Target) HostDir() string {
	if t.Scheme == "ssh" && (t.Port == "" || t.Port == "22") {
		return t.Host
	}
	if (t.Scheme == "https" && t.Port == "443") || (t.Scheme == "http" && t.Port == "80") {
		return t.Host
	}
	if t.Port != "" {
		return t.Host + "~" + t.Port
	}
	return t.Host
}

func (t Target) WebURL() (string, error) {
	if t.Local {
		return "", errors.New("a local-only origin has no website; use an explicit HTTPS target")
	}
	if t.Scheme == "ssh" {
		if t.Host != "github.com" && t.Host != "gitlab.com" && t.Host != "bitbucket.org" {
			return "", errors.New("cannot map SSH host/alias to a website; use an explicit HTTPS target")
		}
		if t.Port != "" && t.Port != "22" {
			return "", errors.New("cannot map this SSH port to a website; use an explicit HTTPS target")
		}
		return (&url.URL{Scheme: "https", Host: t.Host, Path: "/" + t.Owner + "/" + t.Repo}).String(), nil
	}
	return (&url.URL{Scheme: t.Scheme, Host: t.HostPort(), Path: "/" + t.Owner + "/" + t.Repo}).String(), nil
}

func ExistingDir(p string) (string, error) {
	info, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("target is not a directory")
	}
	return p, nil
}

// Locate returns disk spellings; missing components use a predictable name.
func (t Target) Locate(storage string) (string, error) {
	if t.Local {
		return t.Path, nil
	}
	p := storage
	names := []string{t.HostDir(), t.Owner, t.Repo}
	for i, name := range names {
		fold := i == 0 || t.Host == "github.com"
		if fold {
			name = strings.ToLower(name)
		}
		entries, err := os.ReadDir(p)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		var matches []string
		for _, entry := range entries {
			if (fold && strings.EqualFold(entry.Name(), name)) || (!fold && entry.Name() == name) {
				matches = append(matches, entry.Name())
			}
		}
		if len(matches) > 1 {
			return "", fmt.Errorf("ambiguous repository directories under %s: %s; choose an explicit local path", p, strings.Join(matches, ", "))
		}
		if len(matches) == 1 {
			name = matches[0]
		}
		p = filepath.Join(p, name)
	}
	return p, nil
}

func (t Target) Directory(storage string) (string, error) {
	p, err := t.Locate(storage)
	if err != nil {
		return "", err
	}
	if !t.Local && t.Subdir != "" {
		p = filepath.Join(p, t.Subdir)
	}
	return ExistingDir(p)
}

// Every component is checked, not just a string prefix of the final path.
func NoSymlinks(p string) error {
	p, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	for {
		info, e := os.Lstat(p)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		if e == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink path component is not allowed for managed writes: %s", p)
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return nil
}

var urls = regexp.MustCompile(`(?:https?|ssh)://[^\s'"<>]+`)

func RedactText(s string) string {
	return urls.ReplaceAllStringFunc(s, func(raw string) string {
		u, err := url.Parse(raw)
		if err != nil {
			return "[redacted URL]"
		}
		u.User = nil
		u.RawQuery = ""
		u.Fragment = ""
		return u.String()
	})
}
