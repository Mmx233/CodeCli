# CodeCli

[![License](https://img.shields.io/github/license/Mmx233/CodeCli)](https://github.com/Mmx233/CodeCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/CodeCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/CodeCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/CodeCli)](https://goreportcard.com/report/github.com/Mmx233/CodeCli)

Inspired by Visual Studio Code 'code' command line tool，Implement command-line automation for cloning, automatically determining project type, and launching Jetbrains idea.

Now testing and optimization are only performed on Windows. 

To ensure proper functionality, either select the option to add idea to the PATH during installation or manually add the 'bin' directory of idea to the PATH. In case of an empty project path or when no matching rules are found, an error will be reported, prompting manual execution if no default idea in config.

After downloading the executable file, place it in the 'bin' directory. For the initial usage, run it once to generate the configuration files.

## Config

The configuration file is located at `$HOME/.CodeCli.yaml`, and all fields are mandatory except for `default.idea`.

```yaml
default:
    gitSite: github.com
    username: "" # default username for gitSite
    cmdProgram: "powershell" # program for 'cmd' sub-command
    idea: "" # default idea, such as goland, webstorm, studio
storage:
    projectDir: "" # the storage path for project files
rules: # match from top to bottom
  - idea: webstorm
    file: [package.json]
  - idea: goland
    file: [go.mod]
  - idea: rustrover
    file: [Cargo.toml]
  - idea: pycharm
    file: [pyproject.toml, requirements.txt]
  - idea: studio
    file: [android\build.gradle, build.gradle]
  - idea: idea
    file: [gradlew]
```

## Usage

Using `default.gitSite=github.com` `default.username=Mmx233` as an example. 

### Open project

```shell
# open idea for current directory
code
code .

code github.com/Mmx233/CodeCli # open the specified project
# commands with similar functionality
code https://github.com/Mmx233/CodeCli
code Mmx233/CodeCli
code CodeCli

code CodeCli --idea webstorm # specify to open the project with WebStorm

# open with relative path
code ..
code ./some_dir
code ../some_dir

# open sub dir as sub project
# sub project only work on full addr
code github.com/Mmx233/CodeCli/sub_dir
code github.com/Mmx233/CodeCli/sub_dir/another_subdir
```

### Clear projects

By default, clean projects that haven't been modified in the last 60 days. Projects with uncommitted code changes will be skipped by default.

```shell
code clear
code clear -t 1000h # specify the idle time during cleanup
code clear Mmx233/CodeCli AnotherCodeCli # delete the specified repository

code clear -y # skip the deletion confirmation
code clear -f # forcefully delete even if there are uncommitted changes or it's not a Git directory
```

### Open a command window for project

```shell
code cmd github.com/Mmx233/CodeCli
code cmd CodeCli
```

### Open project address in browser

Only supported on windows.

```shell
code browser Mmx233/CodeCli
```

### Configure CodeCli

```shell
code config list # list all configs

code config default.username=Mmx233 # modify config
# commands with similar functionality
code config default username Mmx233
code config set default.username=Mmx233

code config unset default.username # unset config
```

### Complete Usage

```shell
~$ code --help-long
usage: code [<flags>] <command> [<args> ...]

A project manager command line tool.

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and
                   --help-man).
  -v, --version    Show application version.
      --idea=IDEA  Specify an idea.

Args:
  <addr>  Project addr.

Commands:
  help [<command>...]
    Show help.


  project [<flags>] [<addr>]
    Open projects.

    --idea=IDEA  Specify an idea.

  clear [<flags>] [<addr>...]
    Auto clear outdated projects.

    -t, --time=1440h  Clean up projects that have not been used for how long.
    -y, --yes         Confirm delete.
    -f, --force       Force delete, skip confirm and checks.

  cmd [<addr>]
    Open project terminal.


  config list
    List all configs.


  config set <field> [<value>]
    Set config.


  config unset <field>
    Clear config.


  browser [<addr>]
    Open project in browser.
```
