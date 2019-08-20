package stringer

import (
	"fmt"
	"reflect"
)

func Stringify(v interface{}) string {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Ptr:
		return fmt.Sprintf("%+v", value.Elem())
	default:
		return fmt.Sprintf("%+v", v)
	}
}
