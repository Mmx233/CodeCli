package models

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

type Commands struct {
	App     *kingpin.Application
	Project ProjectCommand
	Clear   ClearCommand
}

type ProjectCommand struct {
	*kingpin.CmdClause
	Addr string
}

type ClearCommand struct {
	*kingpin.CmdClause
	Duration *time.Duration
}
