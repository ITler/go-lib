package misc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

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
// Deletion is done by explicitly setting a field's empty value.
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
