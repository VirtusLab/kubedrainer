package stringer

import (
	"fmt"
	"reflect"
)

// Stringify converts a given interface into a human readable string
func Stringify(v interface{}) string {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Ptr:
		return fmt.Sprintf("%+v", value.Elem())
	default:
		return fmt.Sprintf("%+v", v)
	}
}
