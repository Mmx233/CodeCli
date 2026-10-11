# CodeCli

[![License](https://img.shields.io/github/license/Mmx233/CodeCli)](https://github.com/Mmx233/CodeCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/CodeCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/CodeCli/releases)

Clone projects and open them in an IDE, terminal or browser. The executable is **code** (`code.exe` on Windows).

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

`code <target>` is shorthand for `code open <target>`; omit the target to use the current directory. `--web`, `--terminal` and `--ide` are mutually exclusive.

Clone reuses existing checkouts by default. Its update flags are mutually exclusive; `--pull` only fast-forwards, and pull/rebase require a clean branch with an upstream. Conflicts remain for manual resolution.

`clone` and `pwd` print an absolute path. Errors return nonzero; Git output goes to stderr.

## Targets and storage

Configure `default.username` to use a bare repository name. Examples:

| Input                                             | Meaning                                |
|---------------------------------------------------|----------------------------------------|
| `CodeCli`                                         | default Git host/default owner/CodeCli |
| `Mmx233/CodeCli`                                  | default Git host/Mmx233/CodeCli        |
| `github.com/Mmx233/CodeCli`                       | explicit host/owner/repo               |
| `Mmx233/CodeCli/internal`                         | GitHub repository subdirectory         |
| `https://github.com/Mmx233/CodeCli.git`           | HTTPS remote                           |
| `ssh://git@github.com/Mmx233/CodeCli.git`         | SSH remote                             |
| `git@github.com:Mmx233/CodeCli.git`               | SCP-style SSH remote                   |
| `.`, `..`, `./folder`, `../folder`, absolute path | local directory                        |

Use `./name` for a local directory; bare names always mean repository shorthand. Windows drive and UNC paths are also supported. Remote subdirectories are currently supported only for GitHub, and `--web` opens the repository homepage.

Projects are stored under `<projectDir>/<host>/<owner>/<repo>`. GitHub identities ignore case: existing directories are reused, new names are lowercase, and duplicate case variants produce an error. Other hosts and local paths retain their own case rules.

Clone never overwrites an existing destination or deletes the final checkout on failure.

## Configuration

Configuration lives in `~/.CodeCli.yaml`. A missing file uses defaults; only explicit configuration changes write it.

```sh
code config get
code config get default.username
code config set default.username Mmx233
code config set default.idea goland
code config set storage.projectDir /absolute/path/to/projects
code config unset default.idea
```

`set` changes scalar keys shown below; `unset` restores defaults. Edit `rules` in YAML. The default project directory is `~/project`.

```yaml
default:
  gitSite: github.com
  username: ""
  cmdProgram: "" # empty selects the platform default
  idea: ""
storage:
  projectDir: ~/project
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

IDE selection: `--ide`, then the first matching rule, then `default.idea`.

## Desktop platforms

| Platform | IDE | Terminal default | Browser |
| --- | --- | --- | --- |
| macOS | PATH launcher or JetBrains .app | Terminal | `open` |
| Linux | PATH launcher | xdg-terminal-exec or common terminal emulator | `xdg-open` |
| Windows | PATH launcher or JetBrains `64.exe` | PowerShell | `rundll32` |

Set `default.cmdProgram` to a terminal executable, or leave it empty for the platform default. On macOS/Linux, shells such as `zsh` and `bash` cannot create a desktop terminal window by themselves.

Run GUI commands from a desktop session. Use `code pwd` to work in the current shell instead.

## Pruning

```sh
code prune Mmx233/CodeCli                  # preview only
code prune Mmx233/CodeCli --apply          # actually remove
code prune Mmx233/CodeCli --force          # preview despite local data
code prune Mmx233/CodeCli --force --apply  # knowingly discard protected local data
```

Prune takes explicit repository targets and previews by default. Add `--apply` to delete them.

Local changes, untracked/ignored files, stash, local-only branches/commits and tags are protected. `--force` permits discarding that data, but always retains managed-path, symlink and worktree checks. Inspect the preview before applying it.
