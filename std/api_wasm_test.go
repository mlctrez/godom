package std

import (
	"github.com/mlctrez/godom"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEach(t *testing.T) {
	a := require.New(t)
	m := godom.Global().Get("Map").New()
	m.Call("set", "a", "av")
	m.Call("set", "b", "bv")
	var result []string
	Each(m.Call("entries"), func(val godom.Value) {
		result = append(result, val.Index(0).String())
	})
	a.Equal(2, len(result))
	a.Equal("a", result[0])
	a.Equal("b", result[1])
}

func TestMapEach(t *testing.T) {
	a := require.New(t)
	m := godom.Global().Get("Map").New()
	m.Call("set", "a", "av")
	m.Call("set", "b", "bv")
	result := map[string]string{}

	MapEach(m.Call("entries"), func(k, v godom.Value) {
		result[k.String()] = v.String()
	})
	a.Equal(2, len(result))
	a.Equal("av", result["a"])
	a.Equal("bv", result["b"])
}

func TestObject(t *testing.T) {
	r := require.New(t)
	// make sure cache is primed, probably is from other tests
	Object()
	r.NotNil(cachedObject)
	r.True(cachedObject.Equal(Object()))
}

func TestOwnPropertyNames(t *testing.T) {
	r := require.New(t)
	r.NotNil(r)
	g := godom.Global()

	propNames := OwnPropertyNames(Proto(g.Get("Array")))
	r.Contains(propNames, "length")

	propNames = OwnPropertyNames(Proto(g.Get("Map")))
	r.Contains(propNames, "length")
	r.NotContains(propNames, "entries")

	propNames = OwnPropertyNames(Proto(g.Get("Map").New()))
	r.NotContains(propNames, "length")
	r.Contains(propNames, "entries")
}
