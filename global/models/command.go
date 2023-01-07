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
}

type ProjectCommand struct {
	*kingpin.CmdClause
	Addr string
	Idea string
}

type ClearCommand struct {
	*kingpin.CmdClause
	Duration time.Duration
	Yes      bool
}

type CmdCommand struct {
	*kingpin.CmdClause
	Addr string
}
