package misc_test

import (
	"testing"

	"github.com/signavio/core-bootstrap-sap/lib/misc"

	"github.com/stretchr/testify/assert"
)

func TestStringJoinWithoutEmpties(t *testing.T) {
	tests := []struct {
		value          []string
		sep            string
		want           string
		wantKeepSpaces string
	}{
		{
			value:          []string{" foo ", " bar "},
			sep:            " - ",
			want:           "foo - bar",
			wantKeepSpaces: "foo - bar",
		},
		{
			value:          []string{" foo", " bar"},
			sep:            ":",
			want:           "foo:bar",
			wantKeepSpaces: "foo:bar",
		},
		{
			value:          []string{"foo", "", " bar "},
			sep:            "-",
			want:           "foo-bar",
			wantKeepSpaces: "foo-bar",
		},
		{
			value:          []string{"foo ", " ", " bar", " "},
			sep:            " - ",
			want:           "foo - bar",
			wantKeepSpaces: "foo -  - bar -",
		},
		{
			value:          []string{"foo", "", " bar", " "},
			sep:            " - ",
			want:           "foo - bar",
			wantKeepSpaces: "foo - bar -",
		},
		{
			value:          []string{" foo ", "", " bar ", ""},
			sep:            "",
			want:           "foo bar",
			wantKeepSpaces: "foo bar",
		},
		{
			value:          []string{"foo", " ", "bar", " "},
			sep:            "",
			want:           "foo bar",
			wantKeepSpaces: "foo  bar",
		},
	}
	for num, tt := range tests {
		assert.Equal(t, tt.want, misc.StringJoinWithoutEmpties(tt.value, tt.sep, false),
			"test case %d/%d - input: %#v", num, len(tests)-1, tt.value)
		assert.Equalf(t, tt.wantKeepSpaces, misc.StringJoinWithoutEmpties(tt.value, tt.sep, true),
			"test case %d/%d - input: %#v", num, len(tests)-1, tt.value)
	}
}

func TestEnrichMapStringInterface(t *testing.T) {
	base := map[string]interface{}{
		"update": "base",
		"keep":   "me",
		"sub": map[string]interface{}{
			"update": "base",
			"keep":   "me",
		},
		"delete_map": map[string]interface{}{
			"val": "val",
		},
		"delete_nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"val":  "val",
				"val2": 5,
			},
		},
		"delete_primitive_int":    1,
		"delete_primitive_string": "to delete",
		"delete_only_struct": struct{ s string }{
			s: "delete whole struct",
		},
		"keep_nil": nil,
	}
	additions := map[string]interface{}{
		"add":    "new",
		"update": "update",
		"sub": map[string]interface{}{
			"add":    "new",
			"update": "update",
		},
		"new": map[string]interface{}{
			"foo": "foo",
		},
		"delete_map": map[string]interface{}{},
		"delete_nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"val":  "",
				"val2": 0,
			},
		},
		"delete_primitive_int":    0,
		"delete_primitive_string": "",
		"delete_only_struct":      struct{ s string }{},
		"array":                   []int{0, 1, 2},
	}
	expected := map[string]interface{}{
		"add":      "new",
		"update":   "update",
		"keep":     "me",
		"keep_nil": nil,
		"sub": map[string]interface{}{
			"add":    "new",
			"update": "update",
			"keep":   "me",
		},
		"new": map[string]interface{}{
			"foo": "foo",
		},
		"array": []int{0, 1, 2},
	}
	actual, err := misc.FlatMergeMapStringInterface(base, additions)
	assert.NoError(t, err, "Should only fail on missing implementation")
	assert.Equal(t, expected, actual, "Should match expectation")

	additions["unsupported_structt"] = struct{ s string }{
		s: "fail",
	}
	actual, err = misc.FlatMergeMapStringInterface(base, additions)
	assert.Errorf(t, err, "Should fail on unsupported implementation: %#v", actual)
}

func TestStructToMapStringInterface(t *testing.T) {
	testStruct := struct {
		Foo string
		Sub struct {
			Bar string
		}
		unconsidered string
	}{
		Foo: "foo",
		Sub: struct{ Bar string }{
			Bar: "bar",
		},
		unconsidered: "unexposed",
	}

	expected := map[string]interface{}{
		"Foo": "foo",
		"Sub": map[string]interface{}{
			"Bar": "bar",
		},
	}

	actual := misc.StructToMapStringInterface(testStruct)
	assert.Equalf(t, expected, actual,
		"Returned map should match expectation, but got: %#v", actual)
}

func TestUpcaseFirstStrict(t *testing.T) {
	tests := []struct {
		expected string
		in       string
	}{
		{
			expected: "First",
			in:       "first",
		},
		{
			expected: "Second",
			in:       "SecONd",
		},
		{
			expected: "Third",
			in:       "THIRD",
		},
	}

	for num, test := range tests {
		assert.Equalf(t, test.expected, misc.UpcaseFirstStrict(test.in), "Test #%s", num)
	}
}
