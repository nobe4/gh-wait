package ready

import (
	"github.com/nobe4/gh-wait/internal/checker/approved"
	"github.com/nobe4/gh-wait/internal/checker/green"
	"github.com/nobe4/gh-wait/internal/github"
)

type Checker struct{}

func (Checker) Check(p github.Pull) (bool, string, error) {
	a, msgA, errA := approved.Checker{}.Check(p)
	if errA != nil {
		return false, "", errA //nolint:wrapcheck // Allow github's errors
	}

	g, msgG, errG := green.Checker{}.Check(p)
	if errG != nil {
		return false, "", errG //nolint:wrapcheck // Allow github's errors
	}

	return a && g, msgA + "\n" + msgG, nil
}
