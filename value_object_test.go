package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValueObject_New(t *testing.T) {
	a := assert.New(t)
	vo := (*ValueObject)(nil).New()
	a.IsType(&ValueObject{}, vo)
}
