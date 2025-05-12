package green

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"

	"github.com/nobe4/gh-wait/internal/github"
)

const (
	// https://docs.github.com/en/graphql/reference/enums#checkconclusionstate
	// Longest value is 15 characters.
	maxStatusLength = 15

	// https://docs.github.com/en/graphql/reference/enums#checkstatusstate
	// Longest value is 11 characters.
	maxConclusionLength = 11

	// TODO: get terminal width.
	// 80 - 1 - 15 - 1 - 11 = 52.
	maxNameLength = 52
)

type Checker struct{}

type CheckRun struct {
	Name       string
	Status     string
	Conclusion string
}

type CheckRuns []CheckRun

func (Checker) Check(p github.Pull) (bool, string, error) {
	//revive:disable:nested-structs // Allow nested structure for simplicity
	var query struct {
		Repository struct {
			PullRequest struct {
				StatusCheckRollup struct {
					Contexts struct {
						Nodes []struct {
							CheckRun CheckRun `graphql:"...on CheckRun"`
						}
					} `graphql:"contexts(first: 100)"`
				}
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]any{
		"owner":  graphql.String(p.Owner),
		"repo":   graphql.String(p.Repo),
		"number": graphql.Int(p.Number),
	}

	if err := github.GraphQL(&query, variables); err != nil {
		return false, "", err //nolint:wrapcheck // Allow github's errors
	}

	checks := CheckRuns{}

	green := true

	for _, node := range query.Repository.PullRequest.StatusCheckRollup.Contexts.Nodes {
		checks = append(checks, node.CheckRun)

		if node.CheckRun.Status != "COMPLETED" || node.CheckRun.Conclusion != "SUCCESS" {
			green = false
		}
	}

	out := checks.String()

	if green {
		out += "All checks have completed successfully"
	} else {
		out += "Some checks are still needed"
	}

	return green, out, nil
}

//revive:disable:cognitive-complexity // This doesn't need to be simplified.
func (r CheckRuns) String() string {
	nameLength := 0
	statusLength := 0
	conclusionLength := 0

	for _, check := range r {
		if l := len(check.Name); l > nameLength {
			nameLength = l
		}

		if l := len(check.Status); l > statusLength {
			statusLength = l
		}

		if l := len(check.Conclusion); l > conclusionLength {
			conclusionLength = l
		}
	}

	nameLength = min(nameLength, maxNameLength)
	statusLength = min(statusLength, maxStatusLength)
	conclusionLength = min(conclusionLength, maxConclusionLength)

	out := ""

	for _, check := range r {
		out += fmt.Sprintf(
			"%-*s %-*s %-*s\n",
			nameLength, check.Name[:min(nameLength, len(check.Name))],
			statusLength, check.Status,
			conclusionLength, check.Conclusion,
		)
	}

	return out
}
