package md_test

import (
	"fmt"
	"testing"

	"github.com/itler/go-lib/md"
	"github.com/stretchr/testify/assert"
)

func TestMakeLink(t *testing.T) {
	testCases := []struct {
		url     string
		caption string
		want    string
	}{
		{
			want: fmt.Sprintf("%s", md.ExampleURL),
		},
		{
			caption: "example",
			want:    fmt.Sprintf("[%s](%s)", "example", md.ExampleURL),
		},
	}
	for _, tc := range testCases {
		got := md.MakeLink(tc.url, tc.caption)
		assert.Equal(t, tc.want, got)
	}
}
