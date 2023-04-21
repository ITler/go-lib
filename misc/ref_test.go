package misc_test

import (
	"testing"

	"github.com/itler/go-lib/misc"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	var foo string = "foo"
	assert.Equal(t, &foo, misc.Ref(foo))
}

func TestDeref(t *testing.T) {
	var foo *string
	assert.Equal(t, "", misc.Deref(foo))
	foo = misc.Ref("foo")
	assert.Equal(t, "foo", misc.Deref(foo))

	var bar *int
	assert.Equal(t, 0, misc.Deref(bar))
	bar = misc.Ref(5)
	assert.Equal(t, 5, misc.Deref(bar))
}
