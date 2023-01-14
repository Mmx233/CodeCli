package models

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

type Commands struct {
	App     *kingpin.Application
	Project ProjectCommand
	Clear   ClearCommand
	Cmd     CmdCommand
	Config  ConfigCommand
}

type ProjectCommand struct {
	*kingpin.CmdClause
	Addr string
	Idea string
}

type ClearCommand struct {
	*kingpin.CmdClause
	Addresses []string
	Duration  time.Duration
	Yes       bool
	Force     bool
}

type CmdCommand struct {
	*kingpin.CmdClause
	Addr string
}

type ConfigCommand struct {
	*kingpin.CmdClause
	List  ConfigListCommand
	Set   ConfigSetCommand
	Unset ConfigUnsetCommand
}

type ConfigListCommand struct {
	*kingpin.CmdClause
}

type ConfigSetCommand struct {
	*kingpin.CmdClause
	Field string
	Value string
}

type ConfigUnsetCommand struct {
	*kingpin.CmdClause
	Field string
}
