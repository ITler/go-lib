package ease

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/rs/zerolog/log"
)

// ReadFile reads a file by name and returns its content as string
func ReadFile(file string) string {
	content, err := ioutil.ReadFile(file) // the file is inside the local directory
	if err != nil {
		log.Error().Stack().Err(err).Msg("Unable to read file")
	}
	return string(content)
}

// UnmarshalTestValue takes a jsonInput and a proper struct type as value
// and tries to unmarshal the JSON to the struct
func UnmarshalTestValue(t *testing.T, jsonInput string, value interface{}) {
	if err := json.Unmarshal([]byte(jsonInput), &value); err != nil {
		var fnErrMsg func(msg string)
		if t != nil {
			fnErrMsg = func(msg string) {
				t.Fatal(msg)
			}
		} else {
			fnErrMsg = log.Panic().Msg
		}
		fnErrMsg("Input does not seem to be parsable to JSON")
	}
}
