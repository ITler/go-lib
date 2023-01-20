package api_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/ITler/go-lib/api"
	"github.com/ITler/go-lib/misc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Logger()
}

const (
	envError = "Setting a variable during tests must be supported by the environment"
)

type HTTPTestClient struct {
	ctx context.Context
	ts  oauth2.TokenSource
}

func (o HTTPTestClient) Token() (*oauth2.Token, error) {
	return nil, errors.New("some fail")
}

func TestNewOAuthClient(t *testing.T) {
	t.Run("happy path not failing with valid token env var defined", func(t *testing.T) {
		for _, e := range api.WellKnownGithubTokenVarNames {
			assert.NoError(t, os.Setenv(e, "123"), envError)
			defer misc.GetAndUnsetEnv(e)
			_, gotErr := api.NewOAuthClient(nil, nil)
			assert.NoError(t, gotErr)
		}

	})
	t.Run("fail when no token env var is defined", func(t *testing.T) {
		for _, e := range api.WellKnownGithubTokenVarNames {
			misc.GetAndUnsetEnv(e)
		}
		_, gotErr := api.NewOAuthClient(nil, nil)
		assert.Error(t, gotErr)
	})

	t.Run("fail when invalid token env var defined", func(t *testing.T) {
		for _, e := range api.WellKnownGithubTokenVarNames {
			assert.NoError(t, os.Setenv(e, ""), envError)
			defer misc.GetAndUnsetEnv(e)
			_, gotErr := api.NewOAuthClient(nil, nil)
			assert.Error(t, gotErr)
		}

	})
	t.Run("fail on invalid envvar token provider", func(t *testing.T) {
		input := api.EnvVarTokenProvider{
			EnvvarNames: []string{},
		}
		_, gotErr := api.NewOAuthClient(nil, &input)
		assert.Error(t, gotErr)
	})
	t.Run("ensure flexible creation procedure for oAuth2 http client", func(t *testing.T) {
		oaccf := func(ctx context.Context, ts oauth2.TokenSource) (api.HTTPClient, error) {
			return &HTTPTestClient{
				ctx: ctx,
				ts:  ts,
			}, nil
		}
		testToken := "123"
		gotRaw, _ := api.NewOAuthClient(nil, api.StringToken(testToken), oaccf)
		got, errC := gotRaw.(*HTTPTestClient)
		assert.True(t, errC)
		token, _ := got.ts.Token()
		assert.NotNil(t, token)
		assert.Equal(t, testToken, token.AccessToken)
		assert.Equal(t, context.Background(), got.ctx)
	})
}

func TestNewOAuthClientDefault(t *testing.T) {
	t.Run("happy path not failing with valid token env var defined", func(t *testing.T) {
		t.Skip()
		for _, e := range api.WellKnownGithubTokenVarNames {
			assert.NoError(t, os.Setenv(e, "123"), envError)
			_, gotErr := api.NewOAuthClientDefault(nil, nil)
			defer misc.GetAndUnsetEnv(e)
			assert.NoError(t, gotErr)
		}

	})
	t.Run("fail when no valid token was provided", func(t *testing.T) {
		_, gotErr := api.NewOAuthClientDefault(nil, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ""}))
		assert.Error(t, gotErr)
	})
	t.Run("fail on failing token source", func(t *testing.T) {
		for _, e := range api.WellKnownGithubTokenVarNames {
			misc.GetAndUnsetEnv(e)
		}
		_, gotErr := api.NewOAuthClientDefault(nil, &HTTPTestClient{})
		assert.Error(t, gotErr)
	})
}
