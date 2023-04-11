package deref_test

import (
	"testing"

	"github.com/itler/go-lib/deref"
	"github.com/stretchr/testify/assert"
)

func TestDerefString(t *testing.T) {
	t.Run("defined input", func(t *testing.T) {
		in := "hello"
		got := deref.String(&in)
		assert.Equal(t, "hello", got)
	})

	t.Run("nil input", func(t *testing.T) {
		var in *string = nil
		got := deref.String(in)
		assert.Equal(t, "", got)
	})
}
