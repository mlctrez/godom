package gdutil

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReflectTypeOf(t *testing.T) {
	a := require.New(t)
	a.Equal(nil, ReflectTypeOf(nil))
	a.Equal("string", fmt.Sprintf("%v", ReflectTypeOf("string")))

}
