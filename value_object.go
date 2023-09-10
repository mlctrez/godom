package godom

type ValueObject struct{ m map[string]interface{} }

func (vo *ValueObject) new() *ValueObject {
	return &ValueObject{m: make(map[string]interface{})}
}

func (vo *ValueObject) New() any            { return (*ValueObject)(nil).new() }
func (vo *ValueObject) Get(p string) any    { return vo.m[p] }
func (vo *ValueObject) Set(p string, x any) { vo.m[p] = x }
func (vo *ValueObject) Delete(p string)     { delete(vo.m, p) }

// IsInstance used by (*value).InstanceOf
func (vo *ValueObject) IsInstance(other Value) any {
	return other.Type() == TypeObject
}
