package callback

import (
	"github.com/mlctrez/godom"
	"reflect"
)

func Reflect(s any) func(e godom.Element, name, data string) {
	m := make(map[string]func(godom.Element))
	elem := reflect.TypeOf(s).Elem()
	for outer := 0; outer < elem.NumField(); outer++ {
		i := outer
		goName := elem.Field(i).Tag.Get("go")
		if goName != "" {
			field := reflect.ValueOf(s).Elem().Field(i)
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
