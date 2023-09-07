package godom

type Func interface {
	Release()
}

func FuncOf(fn func(this Value, args []Value) any) Func {
	return funcOf(fn)
}
