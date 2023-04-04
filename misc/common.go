package misc

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

// DeferFunc conveniently describes a type for a function
// that is intended to be used with a defer statement
type DeferFunc func()

// ArrayDeleteEmpties removes empty elements from a string array
// and returns a new array with the cleaned result.
//
// Setting keepSpaces will not trim elements before checking for deletion
func ArrayDeleteEmpties(s []string, keepSpaces bool) []string {
	var result []string
	for _, str := range s {
		if !keepSpaces {
			str = strings.TrimSpace(str)
		}
		if str != "" {
			result = append(result, strings.TrimSpace(str))
		}
	}
	return result
}

// CreateTmpDir conveniently creates a temporary directory
// and returns the path to the temporary directory along with
// a function reference rmDir for directory removal.
//
// In case of erros creating the temporary directory, err indicates
// an error
func CreateTmpDir(dirName string) (tmpDir string, rmDir DeferFunc, err error) {
	tmpDir, err = os.MkdirTemp("", "*-"+strings.ReplaceAll(dirName, "/", "-"))
	if err != nil {
		return "", nil, fmt.Errorf("Creating temp dir for file download failed: %w", err)
	}
	rmDir = func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Panic().Msgf("Temp dir '%s' cannot be removed: %s", tmpDir, err)
		}
	}

	return
}

// OpenFileWrite opens a file determined by path for writing.
// The file gets created if it not exists
//
// This function utilizes [os.OpenFile] providing
// `os.O_CREATE|os.O_WRONLY|os.O_TRUNC` and [os.ModePerm] as parameters
func OpenFileWrite(filepath string) (*os.File, DeferFunc, error) {
	accessFlags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC

	file, err := os.OpenFile(filepath, accessFlags, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		file.Close()
	}, err
}

// StringJoinWithoutEmpties joins elements of an array after deleting empty elements.
// For deleting empty array elements [ArrayDeleteEmpties] is used
func StringJoinWithoutEmpties(stringss []string, sep string, keepSpaces bool) string {
	if sep == "" {
		sep = " "
	}
	return strings.TrimSpace(strings.Join(ArrayDeleteEmpties(stringss, keepSpaces), sep))
}

// StructToMapStringInterface turns a given struct to a map[string]interface{}
func StructToMapStringInterface(structt interface{}) (result map[string]interface{}) {
	b, _ := json.Marshal(structt)
	_ = json.Unmarshal(b, &result)
	return
}

// FlatMergeMapStringInterface logically takes the base and merge changes onto it.
// Deletion is done by explicitly setting a field to its empty value.
// Nil is kept, so it could be set, explicitly.
// For more sophisticated merging, use a library like https://github.com/imdario/mergo
func FlatMergeMapStringInterface(base, changes map[string]interface{}) (result map[string]interface{}, err error) {
	result = base
	for k, v := range changes {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			if result[k] != nil {
				result[k], err = FlatMergeMapStringInterface(result[k].(map[string]interface{}), v.(map[string]interface{}))
			} else {
				result[k] = v
			}
			for _, del := range []interface{}{result[k], v} {
				if len(del.(map[string]interface{})) == 0 {
					delete(result, k)
					continue
				}
			}
		case reflect.Struct:
			if v != nil && reflect.ValueOf(v).IsZero() {
				delete(result, k)
				continue
			}
			err = fmt.Errorf("Processing value of '%T' is not implemented", v)
			return nil, err
		default:
			if v != nil && reflect.ValueOf(v).IsZero() {
				delete(result, k)
				continue
			}
			result[k] = v
		}
	}
	return
}

// UpcaseFirstStrict takes a string and only sets first letter to upcase
func UpcaseFirstStrict(s string) string {
	return strings.Title(strings.ToLower(s))
}
