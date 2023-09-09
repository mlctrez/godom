package std

import "github.com/mlctrez/godom"

func Each(iter godom.Value, fn func(godom.Value)) {
	for {
		iterResult := iter.Call("next")
		if iterResult.Get("done").Bool() {
			break
		}
		fn(iterResult.Get("value"))
	}
}

func MapEach(iter godom.Value, fn func(k godom.Value, v godom.Value)) {
	Each(iter, func(v godom.Value) {
		fn(v.Index(0), v.Index(1))
	})
}

var cachedObject godom.Value

func Object() godom.Value {
	if cachedObject == nil {
		cachedObject = godom.Global().Get("Object")
	}
	return cachedObject
}

func Proto(v godom.Value) godom.Value {
	return v.Get("__proto__")
}

func OwnPropertyNames(v godom.Value) (result []string) {
	propArray := Object().Call("getOwnPropertyNames", v)
	MapEach(propArray.Call("entries"), func(k, v godom.Value) {
		result = append(result, v.String())
	})
	return
}
