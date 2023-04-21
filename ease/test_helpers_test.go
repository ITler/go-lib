package ease

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type UnmarshalTest struct {
	val string
}

func (t *UnmarshalTest) Validate() error {
	return nil
}

// func TestUnmarshalTestValues(t *testing.T) {
// 	should := `{"val": "foo"}`
// 	var foo JsonValidator
// 	foo = &UnmarshalTest{
// 		val: "test",
// 	}
// 	UnmarshalTestValue(t, should, &foo)

// 	t.Logf("%+v", foo)
// 	// actual, _ := json.Marshal(&foo)
// 	// assert.Equal(t, should, string(actual), "Marshalling should work, properly")
// }

func TestReadFile(t *testing.T) {
	assert.Regexpf(t,
		"module.*",
		strings.TrimSpace(strings.Split(ReadFile("../go.mod"), "\n")[0]),
		"Should be able to read file contents")
	current := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	assert.Equal(t, "", strings.TrimSpace(ReadFile("non-existing.file")),
		"Should return empty string if file does not exist")
	zerolog.SetGlobalLevel(current)
}
