package base

import (
	dom "github.com/mlctrez/godom"
)

type ValueHelper struct {
	Root dom.Value
}

func (v ValueHelper) Get(path ...string) (out dom.Value, ok bool) {
	for _, part := range path {
		next := v.Root.Get(part)
		if next.IsUndefined() {
			return nil, false
		}
		v.Root = next
	}
	if v.Root.IsUndefined() {
		return nil, false
	}
	return v.Root, true
}

func (v ValueHelper) List(path ...string) (out []dom.Value, ok bool) {
	if list, ok := v.Get(path...); ok {
		length := list.Length()
		for i := 0; i < length; i++ {
			out = append(out, list.Index(i))
		}
		return out, true
	}
	return nil, false
}
