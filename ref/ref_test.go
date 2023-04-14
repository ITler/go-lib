package ref_test

import (
	"testing"

	"github.com/itler/go-lib/ref"
	"github.com/stretchr/testify/assert"
)

func TestRefString(t *testing.T) {
	assert.Equal(t, "foo", *ref.String("foo"))
}
