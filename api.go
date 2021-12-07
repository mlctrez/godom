package godom

// Doc is a helper class for working with the Document interface.
type Doc struct{ Doc Document }

// El creates a new element with optional attributes.
func (d Doc) El(tag string, attributes ...*Attribute) Element {
	c := d.Doc.CreateElement(tag)
	for _, a := range attributes {
		c.SetAttribute(a.Name, a.Value)
	}
	return c
}

// At creates a new attribute.
func (d Doc) At(name string, value interface{}) *Attribute {
	return &Attribute{Name: name, Value: value}
}

const IM = "implement me"
