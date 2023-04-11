package oauth2

import (
	"context"
	"errors"
	"fmt"

	"github.com/itler/go-lib/api"
	goauth2 "golang.org/x/oauth2"
)

// ClientCreationFunc defines the structure of a function,
// which is capable for creating an oauth2 client
type ClientCreationFunc func(ctx context.Context, ts goauth2.TokenSource) (api.HTTPClient, error)

// NewClient conveniently creates an oAuth2 client out of
// the provided factory function referred by oaccf [OAuth2ClientCreationFunc] and
// by using a string token provided via stp [StringTokenProvider]
//
// If the provided context is empty, the context will be initialised
// with [context.Background]
//
// If stp is not provided, the routine defaults to trying to retrieve
// tokens from default environment variables determined by [WellKnownGithubTokenVarNames]
//
// Also if oaccf is not provided, the default oAuth2 client creation mechanism
// is triggered via [NewOAuthClientDefault]
func NewClient(ctx context.Context, stp api.StringTokenProvider, oaccf ...ClientCreationFunc) (api.HTTPClient, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if stp == nil {
		return nil, errors.New("No token provider defined")
	}

	if len(oaccf) == 0 || oaccf[0] == nil {
		oaccf = []ClientCreationFunc{NewClientDefault}
	}

	token, err := stp.Parse()
	if err != nil {
		return nil, err
	}

	return oaccf[0](ctx, NewStaticTokenSource(token))
}

// NewClientDefault creates a default oAuth2 client
func NewClientDefault(ctx context.Context, ts goauth2.TokenSource) (api.HTTPClient, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ts == nil {
		return nil, errors.New("No token source provided for oAuth2 client creation")
	}
	if t, err := ts.Token(); err != nil || !t.Valid() {
		if err != nil {
			return nil, fmt.Errorf("No token was provided from token source: %w", err)
		}
		return nil, errors.New("Provided token is not valid")
	}

	return newClientDefault(ctx, ts)
}

// NewStaticTokenSource conveniently creates a [goauth2.StaticTokenSource]
// based on the provided token, which must be non-empty
func NewStaticTokenSource(token string) goauth2.TokenSource {
	if token == "" {
		return nil
	}
	return goauth2.StaticTokenSource(&goauth2.Token{AccessToken: token})
}

func newClientDefault(ctx context.Context, ts goauth2.TokenSource) (api.HTTPClient, error) {
	return goauth2.NewClient(ctx, ts), nil
}
