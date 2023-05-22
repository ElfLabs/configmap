/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package configmap

import (
	"fmt"
	"strings"

	"github.com/mitchellh/copystructure"
)

// NestedFieldNoCopy returns a reference to a nested field.
// Returns false if value is not found and an error if unable
// to traverse obj.
//
// Note: fields passed to this function are treated as keys within the passed
// object; no array/slice syntax is supported.
func NestedFieldNoCopy(obj map[string]interface{}, fields ...string) (interface{}, bool, error) {
	var val interface{} = obj

	for i, field := range fields {
		if val == nil {
			return nil, false, nil
		}
		if m, ok := val.(map[string]interface{}); ok {
			val, ok = m[field]
			if !ok {
				return nil, false, nil
			}
		} else {
			return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected map[string]interface{}", jsonPath(fields[:i+1]), val, val)
		}
	}
	return val, true, nil
}

func jsonPath(fields []string) string {
	return "." + strings.Join(fields, ".")
}

func setNestedFieldNoCopy(obj map[string]interface{}, value interface{}, fields ...string) error {
	m := obj

	for i, field := range fields[:len(fields)-1] {
		if val, ok := m[field]; ok {
			if valMap, ok := val.(map[string]interface{}); ok {
				m = valMap
			} else {
				return fmt.Errorf("value cannot be set because %v is not a map[string]interface{}", jsonPath(fields[:i+1]))
			}
		} else {
			newVal := make(map[string]interface{})
			m[field] = newVal
			m = newVal
		}
	}
	m[fields[len(fields)-1]] = value
	return nil
}

// NestedString returns the string value of a nested field.
// Returns false if value is not found and an error if not a string.
func NestedString(obj map[string]interface{}, fields ...string) (string, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return "", found, err
	}
	s, ok := val.(string)
	if !ok {
		return "", false, fmt.Errorf("%v accessor error: %v is of the type %T, expected string", jsonPath(fields), val, val)
	}
	return s, true, nil
}

// NestedBool returns the bool value of a nested field.
// Returns false if value is not found and an error if not a bool.
func NestedBool(obj map[string]interface{}, fields ...string) (bool, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return false, found, err
	}
	b, ok := val.(bool)
	if !ok {
		return false, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected bool", jsonPath(fields), val, val)
	}
	return b, true, nil
}

// NestedFloat64 returns the float64 value of a nested field.
// Returns false if value is not found and an error if not a float64.
func NestedFloat64(obj map[string]interface{}, fields ...string) (float64, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return 0, found, err
	}
	f, ok := val.(float64)
	if !ok {
		return 0, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected float64", jsonPath(fields), val, val)
	}
	return f, true, nil
}

// NestedInt64 returns the int64 value of a nested field.
// Returns false if value is not found and an error if not an int64.
func NestedInt64(obj map[string]interface{}, fields ...string) (int64, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return 0, found, err
	}
	i, ok := val.(int64)
	if !ok {
		return 0, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected int64", jsonPath(fields), val, val)
	}
	return i, true, nil
}

// NestedStringSlice returns a copy of []string value of a nested field.
// Returns false if value is not found and an error if not a []interface{} or contains non-string items in the slice.
func NestedStringSlice(obj map[string]interface{}, fields ...string) ([]string, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	m, ok := val.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected []interface{}", jsonPath(fields), val, val)
	}
	strSlice := make([]string, 0, len(m))
	for _, v := range m {
		if str, ok := v.(string); ok {
			strSlice = append(strSlice, str)
		} else {
			return nil, false, fmt.Errorf("%v accessor error: contains non-string key in the slice: %v is of the type %T, expected string", jsonPath(fields), v, v)
		}
	}
	return strSlice, true, nil
}

// NestedSlice returns a deep copy of []interface{} value of a nested field.
// Returns false if value is not found and an error if not a []interface{}.
func NestedSlice(obj map[string]interface{}, fields ...string) ([]interface{}, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	v, ok := val.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected []interface{}", jsonPath(fields), val, val)
	}
	return v, true, nil
}

// NestedStringMap returns a copy of map[string]string value of a nested field.
// Returns false if value is not found and an error if not a map[string]interface{} or contains non-string values in the map.
func NestedStringMap(obj map[string]interface{}, fields ...string) (map[string]string, bool, error) {
	m, found, err := nestedMapNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	strMap := make(map[string]string, len(m))
	for k, v := range m {
		if str, ok := v.(string); ok {
			strMap[k] = str
		} else {
			return nil, false, fmt.Errorf("%v accessor error: contains non-string key in the map: %v is of the type %T, expected string", jsonPath(fields), v, v)
		}
	}
	return strMap, true, nil
}

// nestedMapNoCopy returns a map[string]interface{} value of a nested field.
// Returns false if value is not found and an error if not a map[string]interface{}.
func nestedMapNoCopy(obj map[string]interface{}, fields ...string) (map[string]interface{}, bool, error) {
	val, found, err := NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected map[string]interface{}", jsonPath(fields), val, val)
	}
	return m, true, nil
}

// DeepMapCopy returns a deep copy of a map[string]interface{}
func DeepMapCopy(mp map[string]interface{}) map[string]interface{} {
	out, _ := copystructure.Copy(&mp)
	if res, ok := out.(*map[string]interface{}); ok {
		return *res
	}
	return map[string]interface{}{}
}
