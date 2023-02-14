// Package utils contains the common functions and structures for the application
package utils

import (
	"reflect"
)

// InArray Find element in array
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
