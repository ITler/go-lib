package gh

import (
	"errors"
	"net/http"

	"github.com/ITler/go-lib/api"
	"github.com/ITler/go-lib/api/oauth2"
	"github.com/google/go-github/v49/github"
)

var (
	// WellKnownTokenVarNames lists known env var names for Github related tokens
	WellKnownTokenVarNames = []string{"GITHUB_TOKEN", "GH_TOKEN", "NPM_TOKEN"}
)

// ClientCreationFunc defines the structure of a function,
// which is capable for creating an oauth2 client
type ClientCreationFunc func(*http.Client) (*github.Client, error)

// NewClient conveniently creates a client connection to the Github API
// based on an already authenticated http client connection
//
// If the authenticated client is not provided, a new client will be created
// by calling [NewOAuthClient] trying to find tokens in [WellKnownGithubTokenVarNames]
//
// Providing a parameter for gccf would allow creating the github client
// with a customized function
func NewClient(authenticated *http.Client, gccf ...ClientCreationFunc) (*github.Client, error) {
	if authenticated == nil {
		client, err := oauth2.NewClient(nil, &api.EnvVarTokenProvider{
			EnvvarNames: WellKnownTokenVarNames,
		}, nil)
		if err != nil {
			return nil, err
		}
		authenticated = client.(*http.Client)
	}
	if len(gccf) == 0 || gccf[0] == nil {
		gccf = []ClientCreationFunc{NewClientDefault}
	}

	return gccf[0](authenticated)
}

// NewClientDefault provides a client connection to the Github API
// based on an already authenticated http client connection
func NewClientDefault(authenticated *http.Client) (*github.Client, error) {
	if authenticated == nil {
		return nil, errors.New("No authenticated client connection provided")
	}

	return newClientDefault(authenticated), nil
}

func newClientDefault(authenticated *http.Client) *github.Client {
	return github.NewClient(authenticated)
}
