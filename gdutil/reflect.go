package gdutil

import "reflect"

// ReflectTypeOf returns reflect.TypeOf(in) or nil for nil in.
func ReflectTypeOf(in any) (out any) {
	if in != nil {
		out = reflect.TypeOf(in)
	}
	return out
}
