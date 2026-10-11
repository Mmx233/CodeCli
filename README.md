# CodeCli

[![License](https://img.shields.io/github/license/Mmx233/CodeCli)](https://github.com/Mmx233/CodeCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/CodeCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/CodeCli/releases)

Clone projects and open them in an IDE, desktop terminal, or browser. The executable is **code** (code.exe on Windows). This is a breaking CLI upgrade.

## Commands

```sh
code open Mmx233/CodeCli                  # clone if absent, then open IDE
code Mmx233/CodeCli                       # shorthand for open
code open . --ide goland                 # explicit launcher
code open Mmx233/CodeCli --terminal
code open Mmx233/CodeCli --web             # repository homepage; no clone
code clone Mmx233/CodeCli                 # clone or reuse; print checkout path
code clone Mmx233/CodeCli --fetch          # update refs only
code clone Mmx233/CodeCli --pull           # fast-forward-only pull
code clone Mmx233/CodeCli --rebase         # pull with rebase
code clone . --pull                       # update an existing local Git root
code pwd                                  # current directory
code pwd Mmx233/CodeCli                    # existing directory; never clone
# Change the current shell directory only after successful resolution:
dir=$(code pwd Mmx233/CodeCli) && cd "$dir"
code --help
code --version                            # also -v
code completion zsh                       # also bash, fish, powershell
```

`--web`, `--terminal`, and `--ide` are mutually exclusive and belong to `open` or its root shorthand. Clone's three update flags are mutually exclusive. A fresh clone is not pulled again. With no update flag, an existing Git root is reused. Pull and rebase require a clean working tree (including untracked files), an attached branch with an upstream, and no ongoing merge/rebase or other Git operation. Autostash is disabled; conflicts remain available for manual resolution. Git's existing remotes/upstream and authentication settings are respected. No automatic stash, reset, remote rewrite, or push.

`clone` and `pwd` print one absolute path on stdout on success; errors return a nonzero exit code and no path. Git output goes to stderr. Opening an existing ordinary directory is allowed, but clone reuse requires a Git checkout root.

## Targets and storage

Configure `default.username` to use a bare repository name. Examples:

| Input | Meaning |
| --- | --- |
| `CodeCli` | default Git host/default owner/CodeCli |
| `Mmx233/CodeCli` | default Git host/Mmx233/CodeCli |
| `github.com/Mmx233/CodeCli` | explicit host/owner/repo |
| `Mmx233/CodeCli/internal` | GitHub repository subdirectory |
| `https://github.com/Mmx233/CodeCli.git` | HTTPS remote |
| `ssh://git@github.com/Mmx233/CodeCli.git` | SSH remote |
| `git@github.com:Mmx233/CodeCli.git` | SCP-style SSH remote |
| `.`, `..`, `./folder`, `../folder`, absolute path | local directory |

Windows drive and UNC absolute paths use native Windows rules. Bare words always mean repository shorthand, even if a same-named local directory exists; use `./name` for a local directory. A host without a dot needs an explicit URL, such as `https://internal/Team/Repo`. IPv6 hosts, dangerous/unsupported protocols, traversal segments and ambiguous nested namespaces are rejected. Remote subdirectory syntax is initially supported only for GitHub; a browser opens the repository homepage. Clone requires the repository root.

Managed storage is `<projectDir>/<host>/<owner>/<repo>`. Only GitHub owner/repository identities are case insensitive. Existing directories keep their real spelling; new GitHub directories use lowercase. Multiple case variants produce a conflict error, even when one exactly matches the input. Resolve duplicates manually or use an explicit local path. Other hosts retain namespace case. Local/subdirectory names follow the filesystem's own case semantics. A nondefault port is stored as `host~port`.

Clone writes to an owned staging directory and publishes without replacing an existing destination. Failed clones never recursively remove the final checkout. Managed clone/prune operations share a storage lock; a simultaneous operation may ask you to retry. After a crashed process, remove `<projectDir>/.codecli.lock` only after confirming no clone/prune is running. Do not move managed directories or run other Git mutations concurrently with pruning.

## Configuration

`~/.CodeCli.yaml` is read lazily. A missing file uses in-memory defaults; help, version, completion and read-only commands do not create it. Only explicit configuration changes write it. Existing keys `idea` and `cmdProgram` are preserved.

```sh
code config get
code config get default.username
code config set default.username Mmx233
code config set default.idea goland
code config set storage.projectDir /absolute/path/to/projects
code config unset default.idea
```

Unset removes an override and restores its default. Scalar keys are `default.gitSite`, `default.username`, `default.cmdProgram`, `default.idea`, and `storage.projectDir`. Values containing spaces, equals signs or percent signs are ordinary strings. Rules are edited manually. Updates preserve non-target keys, unknown fields and comments using YAML nodes and temporary-file replacement. Unsupported YAML merges/aliases in known configuration fields fail clearly rather than silently changing their meaning. `config get` shows effective known settings; it does not rewrite the file or discard unknown fields.

Defaults (projectDir is `<home>/project`):

```yaml
default:
  gitSite: github.com
  username: ""
  cmdProgram: "" # empty selects the platform default
  idea: ""
storage:
  projectDir: /absolute/path/to/projects
rules:
  - idea: webstorm
    file: [package.json]
  - idea: goland
    file: [go.mod]
  - idea: rustrover
    file: [Cargo.toml]
  - idea: pycharm
    file: [pyproject.toml, requirements.txt]
  - idea: studio
    file: [android/build.gradle, build.gradle]
  - idea: idea
    file: [gradlew]
```

IDE selection: explicit `--ide`, first matching rule, configured default. Generic and empty directories work with an explicit/default IDE. Known launchers are checked before clone where possible; automatic rule detection requires the checkout. A later launch failure preserves a successful clone.

## Desktop platforms

| Platform | IDE | Terminal default | Browser |
| --- | --- | --- | --- |
| macOS | PATH executable, explicit .app path, common JetBrains apps under /Applications or ~/Applications | Terminal via `open -a` | `open` |
| Linux Desktop | PATH executable | first available xdg-terminal-exec, gnome-terminal, konsole, xfce4-terminal, xterm | `xdg-open` |
| Windows | PATH executable, with JetBrains `64.exe` fallback | PowerShell in a new console; `wt` also supported | `rundll32` |

`default.cmdProgram` accepts one executable/path, not a shell command with embedded arguments. Unknown terminal launchers inherit the project as cwd. On macOS/Linux, shell names such as zsh or bash are rejected because starting them alone cannot create a desktop window; unset default.cmdProgram to use the platform default. Existing values are preserved, not silently rewritten. Linux requires DISPLAY or WAYLAND_DISPLAY; macOS SSH sessions return a desktop-session error. Use `pwd` for the current terminal. Headless web requests print a credential-free URL to stderr and still fail. HTTP(S) website mapping preserves ports and removes credentials/query parameters; SSH mapping only recognizes standard-port github.com/gitlab.com/bitbucket.org, never guesses SSH aliases.

A synchronous system launcher returning successfully confirms dispatch. For a long-running IDE/terminal executable, success confirms process creation only; later application failures and GUI readiness cannot be inferred. Linux/Windows desktop behavior still needs native environment acceptance; platform plans and cross-builds alone are insufficient.

## Pruning

```sh
code prune Mmx233/CodeCli                  # preview only
code prune Mmx233/CodeCli --apply          # actually remove
code prune Mmx233/CodeCli --force          # preview despite local data
code prune Mmx233/CodeCli --force --apply  # knowingly discard protected local data
```

Only explicit targets are accepted; there is no age/mtime-based auto-delete. Default protection covers modified/untracked/ignored files, stash, locally unique commits including reflogs, local branch names absent from remote-tracking branches, local tags, and other local refs. Tags are conservatively protected because offline checks cannot prove publication. Remote-tracking refs describe the last fetched state, not current remote availability.

`--force` bypasses these local-data checks, **never** the path boundary. Prune rejects the storage root, ancestors, subdirectories, ordinary folders, ambiguous identities, symlink ancestors, linked worktrees/external gitdirs, and repositories with additional worktrees. The entire target batch is validated first; OS deletion failures are reported and may leave a partially removed batch. Always inspect the preview and stop other processes using the checkouts before applying it.

## Migration

| Previous usage | Replacement |
| --- | --- |
| `code cmd X` | `code open X --terminal` |
| `code browser X` | `code open X --web` |
| `code project X` | `code open X` |
| `code X --idea goland` | `code open X --ide goland` |
| `code clear ...` | explicit `code prune X`, then `--apply` |
| `code config list` | `code config get` |
| implicit configuration setting | `code config set <key> <value>` |

Removed bare command words report migration guidance instead of cloning. To open a repository with such a name, use `code open cmd` (and configure a default owner). No legacy aliases are installed.

## Development and packages

Requires Go 1.25.5 or newer and Git 2.31 or newer.

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o build/code -ldflags '-X main.Version=dev' .
sh scripts/build.sh v2.0.0
# Windows:
# scripts\build.bat v2.0.0
```

The root main.go and cmd package follow Cobra's basic project layout. Package scripts use Go and tar, produce checksums and archives under `build/<version>`, and include Darwin amd64/arm64 plus Linux/Windows amd64/arm64/386. Archives contain `code` or `code.exe`. Windows uses PowerShell and tar.exe. Building does not publish or install anything. `go install` names the executable after the module directory (`CodeCli`); use the build command or release binary when you need `code`.

This command name conflicts with VS Code's launcher. Choose which executable is on PATH; future tap packaging must make that choice explicit. Release publication, Homebrew integration, and installation tests against the published online version are separate steps.
