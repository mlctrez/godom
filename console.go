package godom

type ConsoleApi interface {
	Log(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Error(args ...any)
}

func Console() ConsoleApi {
	return &console{val: Global().Get("console")}
}

var _ ConsoleApi = (*console)(nil)

type console struct {
	val Value
}

func (c *console) Log(args ...any)   { c.val.Call("log", args...) }
func (c *console) Debug(args ...any) { c.val.Call("debug", args...) }
func (c *console) Info(args ...any)  { c.val.Call("info", args...) }
func (c *console) Error(args ...any) { c.val.Call("error", args...) }
