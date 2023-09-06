package convert

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToInt(t *testing.T) {
	as := assert.New(t)
	as.True(ToInt(0) == 0, "assertion failed int")
	as.True(ToInt(int8(0)) == 0, "assertion failed int8")
	as.True(ToInt(int16(0)) == 0, "assertion failed int16")
	as.True(ToInt(int32(0)) == 0, "assertion failed int32")
	as.True(ToInt(int64(0)) == 0, "assertion failed int64")
	as.True(ToInt(float32(0)) == 0, "assertion failed float32")
	as.True(ToInt(float64(0)) == 0, "assertion failed float64")
}

func TestToIntFail(t *testing.T) {
	req := require.New(t)
	defer func() {
		r := recover()
		switch rt := r.(type) {
		case error:
			req.Equal("ToInt failed for bad", rt.Error())
		default:
			req.Fail("recover did not produce error")
		}
	}()
	ToInt("bad")
}

func TestToFloat(t *testing.T) {
	as := assert.New(t)
	as.True(ToFloat(0) == 0, "assertion failed int")
	as.True(ToFloat(int8(0)) == 0, "assertion failed int8")
	as.True(ToFloat(int16(0)) == 0, "assertion failed int16")
	as.True(ToFloat(int32(0)) == 0, "assertion failed int32")
	as.True(ToFloat(int64(0)) == 0, "assertion failed int64")
	as.True(ToFloat(float32(0)) == 0, "assertion failed float32")
	as.True(ToFloat(float64(0)) == 0, "assertion failed float64")
}

func TestToFloatFail(t *testing.T) {
	req := require.New(t)
	defer func() {
		r := recover()
		switch rt := r.(type) {
		case error:
			req.Equal("ToFloat failed for bad", rt.Error())
		default:
			req.Fail("recover did not produce error")
		}
	}()
	ToFloat("bad")
}

func TestStringsAny(t *testing.T) {
	a := assert.New(t)
	stringsAny := StringsAny("1", "2", "3")
	a.True(len(stringsAny) == 3)
	a.Equal("1", stringsAny[0])
}

func TestToReflectArgs(t *testing.T) {
	a := assert.New(t)
	result := FromReflectArgs(ToReflectArgs("a", "b", "c"))
	a.True(len(result) == 3)
	a.Equal("a", result[0])
}
