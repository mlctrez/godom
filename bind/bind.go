package bind

import (
	"github.com/mlctrez/godom"
	"reflect"
)

func Reflect(s any) func(e godom.Element, name, data string) {
	return func(eIn godom.Element, nameIn, dataIn string) {
		elem := reflect.TypeOf(s).Elem()
		for i := 0; i < elem.NumField(); i++ {
			if elem.Field(i).Tag.Get("go") == dataIn {
				field := reflect.ValueOf(s).Elem().Field(i)
				field.Set(reflect.ValueOf(eIn))
			}
		}
	}
}

func Mapper(m map[string]func(godom.Element)) func(e godom.Element, name, data string) {
	return func(ei godom.Element, name string, data string) {
		if fn, ok := m[data]; ok {
			fn(ei)
		}
	}
}
