package gh_test

import (
	"os"
	"testing"

	"github.com/ITler/go-lib/api/gh"
	"github.com/ITler/go-lib/misc"
	"github.com/google/go-github/v49/github"
	"github.com/stretchr/testify/assert"
)

const (
	envError = "Setting a variable during tests must be supported by the environment"
)

func TestNewGhClient(t *testing.T) {
	t.Run("fail Github client creation when default oAuth client cannot be created",
		func(t *testing.T) {
			// type Tokens struct {
			// 	key string
			// 	val string
			// }
			// tokens := []Tokens{}
			for _, key := range gh.WellKnownTokenVarNames {
				misc.GetAndUnsetEnv(key)
				// if err != nil {
				// 	continue
				// }
				// tokens = append(tokens, Tokens{key, val})
			}
			_, gotErr := gh.NewClient(nil)
			assert.Error(t, gotErr)
		})
	t.Run("happy path creating default Github client",
		func(t *testing.T) {
			tokenVar := gh.WellKnownTokenVarNames[0]
			assert.NoError(t, os.Setenv(tokenVar, "123"), envError)
			defer misc.GetAndUnsetEnv(tokenVar)
			_, gotErr := gh.NewClient(nil)
			assert.NoError(t, gotErr)

		})
}

func TestGetNewOpts(t *testing.T) {
	var opts *github.SearchOptions

	t.Run("initial call to function", func(t *testing.T) {
		opts = gh.GetNewOpts(nil, nil)
		assert.Equalf(t, gh.PerPageResultsDefault, opts.ListOptions.PerPage,
			"default options should let API return %v results per query", gh.PerPageResultsDefault)
	})

	t.Run("response with nil when no page is left", func(t *testing.T) {
		resp := github.Response{
			LastPage: 0,
		}
		got := gh.GetNewOpts(nil, &resp)
		assert.Nil(t, got)
	})

	t.Run("provide proper search options when there are pages left", func(t *testing.T) {
		resp := github.Response{
			LastPage: 1,
			NextPage: 2,
		}
		got := gh.GetNewOpts(nil, &resp)
		assert.Equal(t, 2, got.Page,
			"should contain the next page value (%v) from the provided response", resp.NextPage)
	})
}

func TestNewGhClientDefault(t *testing.T) {
	t.Run("no authenticated conneciton provided", func(t *testing.T) {
		_, gotErr := gh.NewClientDefault(nil)
		assert.Error(t, gotErr)
	})
}
