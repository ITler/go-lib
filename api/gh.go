package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v48/github"
)

const (
	// PerPageResultsDefault defines the default page size for search query responses
	PerPageResultsDefault = 100
)

var (
	// WellKnownGithubTokenVarNames lists known env var names for Github related tokens
	WellKnownGithubTokenVarNames = []string{"GITHUB_TOKEN", "GH_TOKEN", "NPM_TOKEN"}
)

// GhClientCreationFunc defines the structure of a function,
// which is capable for creating an oauth2 client
type GhClientCreationFunc func(*http.Client) (*github.Client, error)

// Queryable is able to query Github API and returns a data structure
type Queryable interface {
	Query(context.Context, *github.Client) (interface{}, error)
}

// GetNewOpts provides [github.SearchOptions]
// optionally based on existing search options and an API response.
//
// This facilitates setting proper search options for handling paging, properly,
// and can be used in a do-while loop that encapsulates and accumulates API responses.
//
//	allRepos := []*github.Repository{}
//	for opts = GetNewOpts(opts, nil); opts != nil; opts = GetNewOpts(opts, resp) {
//		res, resp, err = gh.Search.Repositories(ctx, query, opts)
//		if err != nil {
//			return nil, err
//		}
//		allRepos = append(allRepos, res.Repositories...)
//	}
//
// When opts and resp are provided, GetNewOpts will return nil if
// there is no page left to query.
func GetNewOpts(opts *github.SearchOptions, resp *github.Response) *github.SearchOptions {
	currentPage := 0
	if opts == nil {
		opts = &github.SearchOptions{
			Sort:      "",
			Order:     "",
			TextMatch: false,
			ListOptions: github.ListOptions{
				Page:    currentPage,
				PerPage: PerPageResultsDefault,
			},
		}
	}
	if resp != nil {
		if resp.LastPage > 0 {
			opts.Page = resp.NextPage
		} else {
			return nil
		}
	}

	return opts
}

// NewGithubClient creates a client connection to the Github API
// based on an already authenticated http client connection
func NewGithubClient(authenticated *http.Client, gccf GhClientCreationFunc) (*github.Client, error) {
	if authenticated == nil {
		client, err := NewOAuthClient(nil, nil, nil)
		if err != nil {
			return nil, err
		}
		authenticated = client.(*http.Client)
	}
	if gccf == nil {
		gccf = NewGithubClientDefault
	}

	return gccf(authenticated)
}

// NewGithubClientDefault provides a client connection to the Github API
// based on an already authenticated http client connection
func NewGithubClientDefault(authenticated *http.Client) (*github.Client, error) {
	if authenticated == nil {
		return nil, errors.New("No authenticated client connection provided")
	}

	return newGithubClientDefault(authenticated), nil
}

func newGithubClientDefault(authenticated *http.Client) *github.Client {
	return github.NewClient(authenticated)
}
