package api_test

import (
	"os"
	"testing"

	"github.com/ITler/go-lib/api"
	"github.com/ITler/go-lib/misc"
	"github.com/google/go-github/v49/github"
	"github.com/stretchr/testify/assert"
)

func TestNewGhClientTest(t *testing.T) {
	t.Run("happy path creating default Github client", func(t *testing.T) {
		tokenVar := api.WellKnownGithubTokenVarNames[0]
		assert.NoError(t, os.Setenv(tokenVar, "123"), envError)
		defer misc.GetAndUnsetEnv(tokenVar)
		_, gotErr := api.NewGithubClient(nil)
		assert.NoError(t, gotErr)

	})
	t.Run("fail Github client creation when oAuth client cannot be created", func(t *testing.T) {
		_, gotErr := api.NewGithubClient(nil)
		assert.Error(t, gotErr)
	})
}

func TestGetNewOpts(t *testing.T) {
	var opts *github.SearchOptions

	t.Run("initial call to function", func(t *testing.T) {
		opts = api.GetNewOpts(nil, nil)
		assert.Equalf(t, api.PerPageResultsDefault, opts.ListOptions.PerPage,
			"default options should let API return %v results per query", api.PerPageResultsDefault)
	})

	t.Run("response with nil when no page is left", func(t *testing.T) {
		resp := github.Response{
			LastPage: 0,
		}
		got := api.GetNewOpts(nil, &resp)
		assert.Nil(t, got)
	})

	t.Run("provide proper search options when there are pages left", func(t *testing.T) {
		resp := github.Response{
			LastPage: 1,
			NextPage: 2,
		}
		got := api.GetNewOpts(nil, &resp)
		assert.Equal(t, 2, got.Page,
			"should contain the next page value (%v) from the provided response", resp.NextPage)
	})
}

func TestNewGhClientDefaultTest(t *testing.T) {
	t.Run("no authenticated conneciton provided", func(t *testing.T) {
		_, gotErr := api.NewGithubClientDefault(nil)
		assert.Error(t, gotErr)
	})
}
