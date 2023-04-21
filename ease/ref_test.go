package ease_test

import (
	"testing"

	"github.com/itler/go-lib/ease"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	var foo string = "foo"
	assert.Equal(t, &foo, ease.Ref(foo))
}

func TestDeref(t *testing.T) {
	var foo *string
	assert.Equal(t, "", ease.Deref(foo))
	foo = ease.Ref("foo")
	assert.Equal(t, "foo", ease.Deref(foo))

	var bar *int
	assert.Equal(t, 0, ease.Deref(bar))
	bar = ease.Ref(5)
	assert.Equal(t, 5, ease.Deref(bar))
}
