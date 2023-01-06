package util

import "errors"

type ErrEmptyDir struct {
	Path string
}

func (a ErrEmptyDir) Error() string {
	return "error empty project dir " + a.Path
}

var (
	ErrEmptyDefaultGitSite  = errors.New("error empty default gitSite config")
	ErrEmptyDefaultUsername = errors.New("error empty default username config")
	ErrUnknownInput         = errors.New("error unknown input")
)
