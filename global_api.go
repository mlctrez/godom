package godom

func Global() Value {
	return global()
}

type DocumentApi interface {
	Node
	CreateElement(tag string) Element
	CreateElementNS(tag string, ns string) Element
	CreateTextNode(text string) Text
	DocumentElement() Element
	Body() Element
	EventListener
	DocApi() Doc
}

type document struct {
	node
}

func (d *document) DocApi() Doc {
	return Doc{Doc: d.This()}
}

func (d *document) CreateElement(tag string) Element {
	return ElementFromValue(d.this.Call("createElement", tag))
}

func (d *document) CreateElementNS(tag string, ns string) Element {
	return ElementFromValue(d.this.Call("createElementNS", tag, ns))
}

func (d *document) CreateTextNode(text string) Text {
	return TextFromValue(d.this.Call("createTextNode", text))
}

func (d *document) DocumentElement() Element {
	// TODO: appropriate caching
	return ElementFromValue(d.this.Get("documentElement"))
}

func (d *document) Body() (body Element) {
	for _, child := range d.DocumentElement().ChildNodes() {
		if child.NodeName() == "body" {
			return child.(*element)
		}
	}
	panic("unable to locate body")
}

func Document() DocumentApi {
	docElement := Global().Get("window").Get("document")
	d := &document{node{this: docElement}}
	d.node.children = d.DocumentElement().ChildNodes()
	return d
}
