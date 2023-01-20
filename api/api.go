// Package api encapsulates convenience functionality for interacting
// with well known APIs and related authentication mechanisms
package api

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

// HTTPClient is used here instead of *http.Client
// Casting can be considered as safe
type HTTPClient interface{}

// OAuth2ClientCreationFunc defines the structure of a function,
// which is capable for creating an oauth2 client
type OAuth2ClientCreationFunc func(ctx context.Context, ts oauth2.TokenSource) (HTTPClient, error)

// NewOAuthClient conveniently creates an oAuth2 client out of
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
func NewOAuthClient(ctx context.Context, stp StringTokenProvider, oaccf ...OAuth2ClientCreationFunc) (HTTPClient, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if stp == nil {
		stp = &EnvVarTokenProvider{
			EnvvarNames: WellKnownGithubTokenVarNames,
		}
	}

	if len(oaccf) == 0 || oaccf[len(oaccf)-1] == nil {
		oaccf = []OAuth2ClientCreationFunc{NewOAuthClientDefault}
	}

	token, err := stp.Parse()
	if err != nil {
		return nil, err
	}

	return oaccf[len(oaccf)-1](ctx, NewOAuthStaticTokenSource(token))
}

// NewOAuthClientDefault creates a default oAuth2 client
func NewOAuthClientDefault(ctx context.Context, ts oauth2.TokenSource) (HTTPClient, error) {
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

	return newOAuthClientDefault(ctx, ts)
}

// NewOAuthStaticTokenSource conveniently creates a [oauth2.StaticTokenSource]
// based on the provided token, which must be non-empty
func NewOAuthStaticTokenSource(token string) oauth2.TokenSource {
	if token == "" {
		return nil
	}
	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
}

func newOAuthClientDefault(ctx context.Context, ts oauth2.TokenSource) (HTTPClient, error) {
	return oauth2.NewClient(ctx, ts), nil
}
