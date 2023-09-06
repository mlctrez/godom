//go:build js && wasm

package godom

import (
	"fmt"
	"testing"
)

func TestDocument_DocumentElement(t *testing.T) {
	fmt.Println("THIS IS THE TEST THAT SHOULD BE WORKING!!!")
	globalClear()
	Global().Document().DocumentElement()
}
