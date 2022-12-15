package misc

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAndUnsetEnv(t *testing.T) {
	key := "test"
	expected := "testvalue"
	os.Setenv(key, expected)
	actual, err := GetAndUnsetEnv(key)
	assert.Equalf(t, expected, actual, "Should return expected value '%s' from env", expected)
	assert.NoError(t, err, "Should not throw error when variable can be unset.")

	key = "never_set"
	expected = ""
	actual, err = GetAndUnsetEnv(key)
	assert.Equal(t, expected, actual, "Should return nothing for unset variable")
	assert.NoError(t, err, "Should not throw error when variable does not exist.")
}
