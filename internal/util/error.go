package util

import "errors"

type ErrUnsupportedProjectOrEmptyDir struct {
	Path string
}

func (a ErrUnsupportedProjectOrEmptyDir) Error() string {
	return "error unsupported project or empty project: " + a.Path
}

var (
	ErrEmptyDefaultGitSite  = errors.New("error empty default gitSite config")
	ErrEmptyDefaultUsername = errors.New("error empty default username config")
	ErrIllegalInput         = errors.New("error illegal input")
	ErrProjectMustDir       = errors.New("error target path not a dir")
)
