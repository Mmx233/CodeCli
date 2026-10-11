package desktop

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Mmx233/CodeCli/internal/config"
)

type Plan struct {
	Program string
	Args    []string
	Dir     string
	Wait    bool
	Console bool
}

func Session(goos string, getenv func(string) string) error {
	if (goos == "linux" && getenv("DISPLAY") == "" && getenv("WAYLAND_DISPLAY") == "") ||
		(goos == "darwin" && (getenv("SSH_CONNECTION") != "" || getenv("SSH_TTY") != "")) {
		return errors.New("no desktop session detected; use code pwd to work in this terminal, or run open from a desktop session")
	}
	return nil
}

func SelectIDE(cfg config.Config, explicit, dir string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	for _, r := range cfg.Rules {
		for _, f := range r.File {
			f = strings.ReplaceAll(f, `\`, "/")
			if !filepath.IsLocal(filepath.FromSlash(f)) {
				return "", errors.New("IDE rule file must be a relative path inside the project")
			}
			if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(f))); err == nil {
				return r.Idea, nil
			} else if !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
		}
	}
	if cfg.Default.Idea != "" {
		return cfg.Default.Idea, nil
	}
	return "", errors.New("no matching IDE; use --ide <launcher> or code config set default.idea <launcher>")
}

func IDE(goos, name, dir string, lookup func(string) (string, error)) (Plan, error) {
	if name == "" {
		return Plan{}, errors.New("IDE launcher is empty")
	}
	if p, err := lookup(name); err == nil {
		return Plan{Program: p, Args: []string{dir}, Dir: dir}, nil
	}
	if goos == "windows" && !strings.ContainsAny(name, `/\`) && filepath.Ext(name) == "" {
		if p, err := lookup(name + "64.exe"); err == nil {
			return Plan{Program: p, Args: []string{dir}, Dir: dir}, nil
		}
	}
	if goos == "darwin" {
		app := name
		if !strings.HasSuffix(app, ".app") {
			known := map[string]string{"goland": "GoLand", "webstorm": "WebStorm", "idea": "IntelliJ IDEA", "rustrover": "RustRover", "pycharm": "PyCharm", "studio": "Android Studio"}
			if n, ok := known[name]; ok {
				app = n
			}
			app += ".app"
		}
		candidates := []string{app}
		if !filepath.IsAbs(app) {
			home, _ := os.UserHomeDir()
			candidates = []string{filepath.Join("/Applications", app), filepath.Join(home, "Applications", app)}
		}
		for _, a := range candidates {
			if info, err := os.Stat(a); err == nil && info.IsDir() {
				p, err := lookup("open")
				if err != nil {
					return Plan{}, err
				}
				return Plan{Program: p, Args: []string{"-a", a, dir}, Dir: dir, Wait: true}, nil
			}
		}
	}
	return Plan{}, fmt.Errorf("IDE launcher %q not found; add it to PATH or configure its executable/app path", name)
}

func Terminal(goos, name, dir string, lookup func(string) (string, error)) (Plan, error) {
	if goos != "windows" {
		switch filepath.Base(name) {
		case "sh", "bash", "zsh", "fish", "dash", "pwsh", "powershell":
			return Plan{}, errors.New("default.cmdProgram names a shell, not a desktop terminal; unset it for the platform default, or configure a terminal launcher; use code pwd for this shell")
		}
	}
	if goos == "darwin" && (name == "" || name == "Terminal" || strings.HasSuffix(name, ".app")) {
		if name == "" {
			name = "Terminal"
		}
		p, err := lookup("open")
		if err != nil {
			return Plan{}, err
		}
		return Plan{Program: p, Args: []string{"-a", name, dir}, Dir: dir, Wait: true}, nil
	}
	if name == "" {
		if goos == "windows" {
			name = "powershell"
		} else if goos == "linux" {
			for _, n := range []string{"xdg-terminal-exec", "gnome-terminal", "konsole", "xfce4-terminal", "xterm"} {
				if _, err := lookup(n); err == nil {
					name = n
					break
				}
			}
		}
	}
	if name == "" {
		return Plan{}, errors.New("no terminal launcher found; set default.cmdProgram or use code pwd")
	}
	p, err := lookup(name)
	if err != nil {
		return Plan{}, fmt.Errorf("terminal launcher %q not found; set default.cmdProgram or use code pwd", name)
	}
	plan := Plan{Program: p, Dir: dir, Console: goos == "windows"}
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(name)), ".exe")
	switch base {
	case "xdg-terminal-exec":
		plan.Args = []string{"--dir=" + dir}
	case "gnome-terminal":
		plan.Args = []string{"--working-directory=" + dir}
		plan.Wait = true
	case "konsole":
		plan.Args = []string{"--workdir", dir}
	case "xfce4-terminal":
		plan.Args = []string{"--working-directory=" + dir}
	case "wt":
		plan.Args = []string{"-d", dir}
	}
	return plan, nil
}

func Web(goos, address string, lookup func(string) (string, error)) (Plan, error) {
	name := ""
	args := []string{address}
	switch goos {
	case "darwin":
		name = "open"
	case "linux":
		name = "xdg-open"
	case "windows":
		name = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", address}
	default:
		return Plan{}, errors.New("unsupported desktop platform")
	}
	p, err := lookup(name)
	if err != nil {
		return Plan{}, fmt.Errorf("browser launcher %s not found", name)
	}
	return Plan{Program: p, Args: args, Wait: true}, nil
}

func (p Plan) Launch(ctx context.Context, stderr io.Writer) error {
	if p.Wait {
		c := exec.CommandContext(ctx, p.Program, p.Args...)
		c.Dir = p.Dir
		c.Stdout = stderr
		c.Stderr = stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("desktop launcher failed: %w", err)
		}
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	c := exec.Command(p.Program, p.Args...)
	c.Dir = p.Dir
	// Detached GUI processes must not keep the parent's terminal pipe open.
	setProcessAttributes(c, p.Console)
	if err := c.Start(); err != nil {
		return fmt.Errorf("start desktop launcher: %w", err)
	}
	go func() { _ = c.Wait() }()
	return nil // OS process creation succeeded; this is not a GUI-ready assertion.
}

func CheckSession() error { return Session(runtime.GOOS, os.Getenv) }
