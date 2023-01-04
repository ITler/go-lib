package conversion_test

import (
	"testing"

	"github.com/ITler/go-lib/conversion"
	"github.com/stretchr/testify/assert"
)

func TestDereferenceString(t *testing.T) {
	t.Run("defined input", func(t *testing.T) {
		in := "hello"
		got := conversion.DereferenceString(&in)
		assert.Equal(t, "hello", got)
	})

	t.Run("nil input", func(t *testing.T) {
		var in *string = nil
		got := conversion.DereferenceString(in)
		assert.Equal(t, "", got)
	})
}
