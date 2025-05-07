package green

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/nobe4/gh-wait/internal/github"
)

type Checker struct{}

type Repo struct {
	Pull Pull `json:"pullRequest"`
}

type Pull struct {
	StatusCheckRollup StatusCheckRollup `json:"statusCheckRollup"`
}

type StatusCheckRollup struct {
	Contexts Contexts `json:"contexts"`
}

type Contexts struct {
	Nodes []CheckRun `json:"nodes"`
}

type CheckRun struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

const query = `
query pullChecks($owner: String!, $repo: String!, $number: Int!){
  repository(owner: $owner, name: $repo){
    pullRequest(number: $number){
      statusCheckRollup {
        contexts(first: 100) {
          nodes{
            ...on CheckRun{
              name,
              status,
              conclusion,
            }
          }
        }
      }
    }
  }
}
`

func (Checker) Check(p github.Pull) (bool, string, error) {
	//nolint:revive // Allow nested structure for simplicity
	var query struct {
		Repository struct {
			PullRequest struct {
				StatusCheckRollup struct {
					Contexts struct {
						Nodes []struct {
							CheckRun struct {
								Name       string
								Status     string
								Conclusion string
							} `graphql:"...on CheckRun"`
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
		return false, "TODO", err
	}

	checks := query.Repository.PullRequest.StatusCheckRollup.Contexts.Nodes

	for _, check := range checks {
		fmt.Printf("%+v\n", check.CheckRun)
	}

	return false, "TODO", nil
}
