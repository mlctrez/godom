package godom

import "testing"

func TestEventTarget_AddEventListener(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("EbentTarget.AddEventListener did not panic")
		}
	}()
	et := &eventTarget{}
	et.AddEventListener("", "")
}
