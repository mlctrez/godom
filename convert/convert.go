package convert

import (
	"fmt"
	"reflect"
)

func ToInt(i interface{}) int {
	switch iv := i.(type) {
	case int:
		return iv
	case int8:
		return int(iv)
	case int16:
		return int(iv)
	case int32:
		return int(iv)
	case int64:
		return int(iv)
	case float32:
		return int(iv)
	case float64:
		return int(iv)
	}
	panic(fmt.Errorf("ToInt failed for %+v", i))
}

func ToFloat(i interface{}) float64 {
	switch iv := i.(type) {
	case int:
		return float64(iv)
	case int8:
		return float64(iv)
	case int16:
		return float64(iv)
	case int32:
		return float64(iv)
	case int64:
		return float64(iv)
	case float32:
		return float64(iv)
	case float64:
		return iv
	}
	panic(fmt.Errorf("ToFloat failed for %+v", i))
}

func ToReflectArgs(args ...interface{}) (result []reflect.Value) {
	result = make([]reflect.Value, len(args))
	for i, arg := range args {
		result[i] = reflect.ValueOf(arg)
	}
	return
}

func FromReflectArgs(args []reflect.Value) (result []interface{}) {
	result = make([]interface{}, len(args))
	for i, arg := range args {
		result[i] = arg.Interface()
	}
	return
}

func StringsAny(args ...string) []any {
	result := make([]any, len(args))
	for i, protocol := range args {
		result[i] = protocol
	}
	return result
}
