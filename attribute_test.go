package godom

import "testing"

func TestAttributeMap_FromMap(t *testing.T) {

}
func TestAttribute_SortByName(t *testing.T) {
	var at Attributes
	at = append(at, &Attribute{Name: "lang", Value: "EN"})
	at = append(at, &Attribute{Name: "data-md-foo", Value: "foo"})
	at.SortByName()
	if at[0].Name != "data-md-foo" {
		t.Fatal("failure to sort attributes by name")
	}
}

func TestAttribute_FromMap(t *testing.T) {
	var am AttributeMap
	// should not panic
	fromMap := am.FromMap()
	if fromMap == nil {
		t.Fatal("from map should not return nil")
	}

	am = AttributeMap{}
	am["one"] = "one"
	am["two"] = "two"

	attributes := am.FromMap()
	if len(attributes) != 2 {
		t.Fatal("length incorrect")
	}
	if attributes[0].Name != "one" {
		t.Fatal("incorrect order")
	}
	if attributes[1].Name != "two" {
		t.Fatal("incorrect order")
	}

}
