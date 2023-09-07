//go:build js && wasm

package godom

import (
	"testing"
)

func TestDocument_DocumentElement(t *testing.T) {
	globalClear()
	Global().Document().DocumentElement()
}
