package callback

import (
	"github.com/mlctrez/godom"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapper(t *testing.T) {
	a := require.New(t)
	var one godom.Element
	mapper := Mapper(map[string]func(godom.Element){
		"one": func(eIn godom.Element) { one = eIn },
	})
	a.NotNil(mapper)
	doc := godom.Document().DocApi()
	a.NotPanics(func() { mapper(doc.El("oneElement"), "name", "one") })
	a.NotPanics(func() { mapper(doc.El("oneElement"), "name", "notMapped") })
	a.NotNil(one)
	a.Equal("oneelement", one.NodeName())
}

type TestStruct struct {
	Button godom.Element `go:"button"`
}

func TestReflect(t *testing.T) {
	a := require.New(t)
	ts := &TestStruct{}
	callBack := Reflect(ts)
	a.NotNil(callBack)
	doc := godom.Document().DocApi()
	a.NotPanics(func() { callBack(doc.El("button"), "name", "button") })
	a.NotPanics(func() { callBack(doc.El("button"), "name", "notMapped") })
	a.NotNil(ts.Button)
	a.Equal("button", ts.Button.NodeName())
}

func TestReflect_panics(t *testing.T) {
	a := require.New(t)

	a.PanicsWithValue("nil ptr", func() { Reflect(nil) })
	a.PanicsWithValue("ptr must be pointer", func() { Reflect(TestStruct{}) })
}
