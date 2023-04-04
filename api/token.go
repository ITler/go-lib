// Package api encapsulates convenience functionality for interacting
// with well known APIs and related authentication mechanisms
package api

import (
	"fmt"

	"github.com/ITler/go-lib/misc"
	"github.com/rs/zerolog/log"
)

// TokenProvider will provide the tokens to be processed by the TokenProvider
type TokenProvider[P any] interface {
	Parse() (P, error)
}

// StringTokenProvider will provide string tokens to be processed by the TokenProvider
type StringTokenProvider TokenProvider[string]

// StringToken is a string that behaves like a [StringTokenProvider]
type StringToken string

// Parse returns the string representation of the underlying type
func (s StringToken) Parse() (string, error) {
	return string(s), nil
}

// EnvVarTokenProvider takes names of env vars, that could hold string tokens
type EnvVarTokenProvider struct {
	EnvvarNames []string
}

// Parse retrieves string tokens from env vars
func (p *EnvVarTokenProvider) Parse() (string, error) {
	return parseEnvvars(p.EnvvarNames)
}

func parseEnvvars(envvarNames []string) (string, error) {
	envvarName, token := "", ""
	for _, candidate := range envvarNames {
		if candidate != "" {
			envvarName = candidate
			token = misc.LookupEnvVar(envvarName)
		}
		if token != "" {
			log.Debug().Msgf("Found token in env var '%s'", envvarName)
			return token, nil
		}
	}

	if envvarName == "" {
		return "", fmt.Errorf("Name of env var to retrieve auth token from was not defined")
	}
	if token == "" {
		log.Debug().Msgf("Cannot retrieve auth token from one of the env var candidates: '%v'", envvarNames)
	}

	return "", nil
}
