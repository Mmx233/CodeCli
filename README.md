# CodeCli

[![Lisense](https://img.shields.io/github/license/Mmx233/CodeCli)](https://github.com/Mmx233/CodeCli/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/CodeCli?color=blueviolet&include_prereleases)](https://github.com/Mmx233/CodeCli/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/CodeCli)](https://goreportcard.com/report/github.com/Mmx233/CodeCli)

受启发于 VsCode code command line，实现 Jetbrain idea 命令行自动克隆、自动判断项目类型、自动启动 idea

目前支持 Goland、Webstorm，需要在安装 idea 时勾选添加到 PATH。空项目路径与不支持的项目类型将报错以等待手动执行

下载可执行文件后放入 bin 中，首次使用需运行一次以生成配置文件

## 配置文件

配置文件位于 `$HOME/.CodeCli.yaml`

```yaml
default:
    gitSite: github.com
    username: "" #默认 clone 用户名
storage:
    projectDir: "" #项目文件存储路径
```

## 使用

以 `default.gitSite=github.com` `default.username=Mmx233` 为例

```shell
code #对当前目录执行打开项目
code github.com/Mmx233/CodeCli #打开指定项目
code Mmx233/CodeCli #同上
code CodeCli #同上
```