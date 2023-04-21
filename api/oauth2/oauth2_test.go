package oauth2_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/itler/go-lib/api"
	"github.com/itler/go-lib/api/gh"
	"github.com/itler/go-lib/api/oauth2"
	"github.com/itler/go-lib/ease"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	goauth2 "golang.org/x/oauth2"
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
	ts  goauth2.TokenSource
}

func (o HTTPTestClient) Token() (*goauth2.Token, error) {
	return nil, errors.New("some fail")
}

func TestLearnGitDownload(t *testing.T) {

	// tmpDir, err := os.MkdirTemp("", "*-"+strings.ReplaceAll(path.Dir(""), "/", "-"))
	// if err != nil {
	// 	t.Logf("Creating temp dir for file download failed: %s", err)
	// }
	// t.Logf("%v", tmpDir)
	// defer func() {
	// 	if err := os.RemoveAll(tmpDir); err != nil {
	// 		t.Logf("Temp dir '%s' cannot be removed: %s", tmpDir, err)
	// 	}
	// }()
	// c, err := api.NewGithubClient(nil)
	// assert.NoError(t, err)
	// u, _, err2 := c.Repositories.GetArchiveLink(context.Background(), "signavio", "cloud-facts-finder", github.Tarball, &github.RepositoryContentGetOptions{
	// 	Ref: "",
	// }, true)
	// t.Logf("%#v", u.String())
	u, err2 := url.Parse("http://foo.com/baz.txt")
	t.Logf("%#v", u.String())
	assert.NoError(t, err2)
}

func TestNewOAuthClient(t *testing.T) {
	validStp := &api.EnvVarTokenProvider{
		EnvvarNames: gh.WellKnownTokenVarNames,
	}
	t.Run("happy path not failing with valid token env var defined", func(t *testing.T) {
		for _, e := range gh.WellKnownTokenVarNames {
			assert.NoError(t, os.Setenv(e, "123"), envError)
			defer ease.GetAndUnsetEnv(e)
			_, gotErr := oauth2.NewClient(nil, validStp)
			assert.NoError(t, gotErr)
		}

	})
	t.Run("fail when no token env var is defined", func(t *testing.T) {
		for _, e := range gh.WellKnownTokenVarNames {
			ease.GetAndUnsetEnv(e)
		}
		_, gotErr := oauth2.NewClient(nil, validStp)
		assert.Error(t, gotErr)
	})

	t.Run("fail when invalid token env var defined", func(t *testing.T) {
		for _, e := range gh.WellKnownTokenVarNames {
			assert.NoError(t, os.Setenv(e, ""), envError)
			defer ease.GetAndUnsetEnv(e)
			_, gotErr := oauth2.NewClient(nil, validStp)
			assert.Error(t, gotErr)
		}

	})

	t.Run("fail on invalid envvar token provider", func(t *testing.T) {
		input := api.EnvVarTokenProvider{
			EnvvarNames: []string{},
		}
		_, gotErr := oauth2.NewClient(nil, &input)
		assert.Error(t, gotErr)
	})

	t.Run("fail on no token provider at all", func(t *testing.T) {
		_, gotErr := oauth2.NewClient(nil, nil)
		assert.Error(t, gotErr)
	})

	t.Run("ensure flexible creation procedure for oAuth2 http client", func(t *testing.T) {
		oaccf := func(ctx context.Context, ts goauth2.TokenSource) (api.HTTPClient, error) {
			return &HTTPTestClient{
				ctx: ctx,
				ts:  ts,
			}, nil
		}
		testToken := "123"
		gotRaw, _ := oauth2.NewClient(nil, api.StringToken(testToken), oaccf)
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
		for _, e := range gh.WellKnownTokenVarNames {
			assert.NoError(t, os.Setenv(e, "123"), envError)
			_, gotErr := oauth2.NewClientDefault(nil, nil)
			defer ease.GetAndUnsetEnv(e)
			assert.NoError(t, gotErr)
		}

	})
	t.Run("fail when no valid token was provided", func(t *testing.T) {
		_, gotErr := oauth2.NewClientDefault(nil, goauth2.StaticTokenSource(
			&goauth2.Token{AccessToken: ""}))
		assert.Error(t, gotErr)
	})
	t.Run("fail on failing token source", func(t *testing.T) {
		for _, e := range gh.WellKnownTokenVarNames {
			ease.GetAndUnsetEnv(e)
		}
		_, gotErr := oauth2.NewClientDefault(nil, &HTTPTestClient{})
		assert.Error(t, gotErr)
	})
}
