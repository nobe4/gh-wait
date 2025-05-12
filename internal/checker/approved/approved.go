package approved

import (
	"fmt"
	"strings"

	"github.com/nobe4/gh-wait/internal/github"
)

type RequestedReviewer struct {
	Login string `json:"login"`
}

type RequestedReviewers []RequestedReviewer

type Pull struct {
	RequestedReviewers RequestedReviewers `json:"requested_reviewers"`
}

type Checker struct{}

func (Checker) Check(p github.Pull) (bool, string, error) {
	result := Pull{}

	if err := github.Get(p.APIString(), &result); err != nil {
		return false, "", err //nolint:wrapcheck // Allow github's errors
	}

	approved := len(result.RequestedReviewers) == 0

	msg := fmt.Sprintf("%s is approved", p)
	if !approved {
		msg = fmt.Sprintf("%s awaits approval from %s", p, result.RequestedReviewers.Logins())
	}

	return approved, msg, nil
}

func (r RequestedReviewers) Logins() string {
	out := []string{}

	for _, v := range r {
		out = append(out, v.Login)
	}

	return strings.Join(out, ", ")
}
