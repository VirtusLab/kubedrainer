package stringer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

// Stringify converts a given interface into a human readable string
func Stringify(v interface{}) string {
	var value reflect.Value
	switch v := v.(type) {
	case nil:
		return "<nil>"
	case reflect.Value:
		value = v
	default:
		value = reflect.ValueOf(v)
	}

	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			return "<nil>"
		}
		return fmt.Sprintf("%+v", Stringify(value.Elem()))
	case reflect.Struct:
		values := make([]string, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			fieldType := value.Type().Field(i)
			field := value.Field(i)
			values[i] = fmt.Sprintf("%s:%s", fieldType.Name, Stringify(field))
		}
		return fmt.Sprintf("{%+v}", strings.Join(values, " "))
	case reflect.Invalid:
		log.Debug().Msgf("invalid value: %s %s %s", reflect.TypeOf(v), reflect.ValueOf(v).Type().Kind(), value.Kind())
		return "?"
	default:
		return fmt.Sprintf("%+v", v)
	}
}
