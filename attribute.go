package godom

import "sort"

type Attribute struct {
	Name  string
	Value interface{}
}

type Attributes []*Attribute
type AttributeMap map[string]interface{}

func (attributes Attributes) SortByName() {
	sortFunc := func(i, j int) bool {
		return attributes[i].Name < attributes[j].Name
	}
	sort.SliceStable(attributes, sortFunc)
}

func (attributeMap AttributeMap) FromMap() Attributes {
	ats := make(Attributes, 0)
	for s, i := range attributeMap {
		ats = append(ats, &Attribute{Name: s, Value: i})
	}
	ats.SortByName()
	return ats
}
