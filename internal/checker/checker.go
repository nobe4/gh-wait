package checker

import (
	"github.com/nobe4/gh-wait/internal/checker/approved"
	"github.com/nobe4/gh-wait/internal/checker/green"
	"github.com/nobe4/gh-wait/internal/checker/merged"
	"github.com/nobe4/gh-wait/internal/github"
)

type Checker interface {
	Check(p github.Pull) (pass bool, msg string, err error)
}

//nolint:ireturn // This return is expected to cover all the implementations.
func Get(name string) Checker {
	return map[string]Checker{
		"merged":   merged.Checker{},
		"approved": approved.Checker{},
		"green":    green.Checker{},
	}[name]
}
