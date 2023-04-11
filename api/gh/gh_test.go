package gh_test

import (
	"os"
	"testing"

	"github.com/itler/go-lib/api/gh"
	"github.com/itler/go-lib/misc"
	"github.com/stretchr/testify/assert"
)

const (
	envError = "Setting a variable during tests must be supported by the environment"
)

func TestNewGhClient(t *testing.T) {
	t.Run("fail Github client creation when default oAuth client cannot be created",
		func(t *testing.T) {
			for _, key := range gh.WellKnownTokenVarNames {
				misc.GetAndUnsetEnv(key)
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

func TestNewGhClientDefault(t *testing.T) {
	t.Run("no authenticated conneciton provided", func(t *testing.T) {
		_, gotErr := gh.NewClientDefault(nil)
		assert.Error(t, gotErr)
	})
}
