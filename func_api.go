package godom

type FuncSignature func(this Value, args []Value) any

type Func interface {
	Release()
}

func FuncOf(fn FuncSignature) Func {
	return funcOf(fn)
}
