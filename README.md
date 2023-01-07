# CodeCli

[![Lisense](https://img.shields.io/github/license/Mmx233/CodeCli)](https://github.com/Mmx233/CodeCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/CodeCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/CodeCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/CodeCli)](https://goreportcard.com/report/github.com/Mmx233/CodeCli)

受启发于 VsCode code command line，实现 Jetbrain idea 命令行自动克隆、自动判断项目类型、自动启动 idea

目前支持 Goland、Webstorm、Android Studio，针对 windows 优化。需要在安装 idea 时勾选添加到 PATH，特别的，Android Studio 需要手动配置 PATH。空项目路径与不支持的项目类型将报错以等待手动执行

下载可执行文件后放入 bin 中，首次使用需运行一次以生成配置文件

## 配置文件

配置文件位于 `$HOME/.CodeCli.yaml`

```yaml
default:
    gitSite: github.com
    username: "" #默认 clone 用户名
    cmdProgram: "powershell" #cmd 指令默认 cmd 程序
    idea: "" #默认 idea，无法判断项目类型时使用。填入 idea 二进制文件名称，如：goland、webstorm、studio
storage:
    projectDir: "" #项目文件存储路径
```

## 使用

以 `default.gitSite=github.com` `default.username=Mmx233` 为例

### 打开项目

```shell
code #对当前目录执行打开项目

code github.com/Mmx233/CodeCli #打开指定项目
#其他相同作用命令
code Mmx233/CodeCli
code CodeCli

code CodeCli --idea webstorm #指定 webstorm 打开该项目
```

### 清理项目

默认清理 60 天未修改的项目，有未提交代码的项目将被跳过

```shell
code clear
code clear -t 1000h #指定闲置时间
```

### 打开目标项目命令行窗口

```shell
code cmd github.com/Mmx233/CodeCli
code cmd CodeCli
```

### 配置 CodeCli

```shell
code config list #列出所有配置

code config default.username=Mmx233 #修改设置
#其他相同作用命令
code config default username Mmx233
code config set default.username=Mmx233

code config unset default.username #清除设置
```

### 全部用法

```shell
usage: code [<flags>] <command> [<args> ...]

A project manager command line tool.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and
                 --help-man).
  -v, --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  project* [<flags>] <addr>
    Open projects.

  clear [<flags>] [<duration>]
    Auto clear outdated projects.

  cmd <addr>
    Open project terminal.

  config list
    List all configs.

  config set* <field> [<value>]
    Set config.
    
  config unset <field>
    Clear config.
```