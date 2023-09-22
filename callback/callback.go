package callback

import (
	"github.com/mlctrez/godom"
	"reflect"
)

func Reflect(ptr any) func(e godom.Element, name, data string) {
	m := make(map[string]func(godom.Element))
	if ptr == nil {
		panic("nil ptr")
	}
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		panic("ptr must be pointer")
	}
	elem := reflect.TypeOf(ptr).Elem()
	for outer := 0; outer < elem.NumField(); outer++ {
		i := outer
		goName := elem.Field(i).Tag.Get("go")
		if goName != "" {
			field := reflect.ValueOf(ptr).Elem().Field(i)
			m[goName] = func(eIn godom.Element) { field.Set(reflect.ValueOf(eIn)) }
		}
	}

	return Mapper(m)
}

func Mapper(m map[string]func(godom.Element)) func(e godom.Element, name, data string) {
	return func(ei godom.Element, name string, data string) {
		if fn, ok := m[data]; ok {
			fn(ei)
		}
	}
}
