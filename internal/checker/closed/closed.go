package closed

import (
	"fmt"

	"github.com/nobe4/gh-wait/internal/github"
)

type Pull struct {
	ClosedAt string `json:"closed_at"`
}

type Checker struct{}

func (Checker) Check(p github.Pull) (bool, string, error) {
	result := Pull{}

	if err := github.Get(p.APIString(), &result); err != nil {
		return false, "", err //nolint:wrapcheck // Allow github's errors
	}

	closed := result.ClosedAt != ""

	msg := fmt.Sprintf("%s is closed", p)
	if !closed {
		msg = fmt.Sprintf("%s is not closed", p)
	}

	return closed, msg, nil
}
