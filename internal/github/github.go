package github

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cli/go-gh/v2/pkg/api"
)

var (
	ErrInvalidPullURL = errors.New("invalid pull request URL")
	ErrCreateClient   = errors.New("failed to create GitHub API client")
	ErrRequest        = errors.New("failed to make request to GitHub API")
)

type Pull struct {
	Owner  string
	Repo   string
	Number int
}

func ParsePull(url string) (Pull, error) {
	const pullParts = 4

	pullRe := regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+)/pull/(\d+)$`)
	matches := pullRe.FindStringSubmatch(url)

	if len(matches) != pullParts {
		return Pull{}, fmt.Errorf("%w: %s", ErrInvalidPullURL, url)
	}

	n, err := strconv.Atoi(matches[3])
	if err != nil {
		return Pull{}, fmt.Errorf("%w: %s", ErrInvalidPullURL, url)
	}

	return Pull{Owner: matches[1], Repo: matches[2], Number: n}, nil
}

func (p Pull) String() string {
	return fmt.Sprintf("%s/%s#%d", p.Owner, p.Repo, p.Number)
}

func (p Pull) APIString() string {
	return fmt.Sprintf("repos/%s/%s/pulls/%d", p.Owner, p.Repo, p.Number)
}

func Get(path string, resp any) error {
	c, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateClient, err)
	}

	if err := c.Get(path, resp); err != nil {
		return fmt.Errorf("%w: %w", ErrRequest, err)
	}

	return nil
}

func GraphQL(q any, variables map[string]any) error {
	c, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateClient, err)
	}

	if err := c.Query("", q, variables); err != nil {
		return fmt.Errorf("%w: %w", ErrRequest, err)
	}

	return nil
}
