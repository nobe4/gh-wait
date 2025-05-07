package merged

import (
	"fmt"

	"github.com/nobe4/gh-wait/internal/github"
)

type Pull struct {
	MergedAt string `json:"merged_at"`
}

type Checker struct{}

func (Checker) Check(p github.Pull) (bool, string, error) {
	result := Pull{}

	if err := github.Get(p.APIString(), &result); err != nil {
		return false, "", err //nolint:wrapcheck // Allow github's errors
	}

	merged := result.MergedAt != ""

	msg := fmt.Sprintf("%s is merged", p)
	if !merged {
		msg = fmt.Sprintf("%s is not merged", p)
	}

	return merged, msg, nil
}
