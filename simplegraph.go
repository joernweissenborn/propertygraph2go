package propertygraph2go

type SimpleVertex struct {
	key        Key
	incoming   []*SimpleEdge
	outgoing   []*SimpleEdge
	properties map[string]interface{}
}

func (v *SimpleVertex) Key() Key { return v.key }

func (v *SimpleVertex) Incoming() (es []Edge) {
	for _, e := range v.incoming {
		es = append(es, e)
	}
	return
}

func (v *SimpleVertex) Outgoing() (es []Edge) {
	for _, e := range v.outgoing {
		es = append(es, e)
	}
	return
}

func (v *SimpleVertex) GetProperty(key string) (property interface{}, err error) {
	property, ok := v.properties[key]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	return
}

func (v *SimpleVertex) SetProperty(key string, property interface{}) (err error) {
	v.properties[key] = property
	return
}

func (v *SimpleVertex) Properties() (properties map[string]interface{}, err error) {
	properties = v.properties
	return
}

type SimpleEdge struct {
	key        Key
	label      string
	head       *SimpleVertex
	tail       *SimpleVertex
	properties map[string]interface{}
}

func (e *SimpleEdge) Key() Key { return e.key }

func (e *SimpleEdge) Label() string { return e.label }

func (e *SimpleEdge) Head() Vertex { return e.head }

func (e *SimpleEdge) Tail() Vertex { return e.tail }

func (e *SimpleEdge) GetProperty(key string) (property interface{}, err error) {
	property, ok := e.properties[key]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	return
}

func (e *SimpleEdge) SetProperty(key string, property interface{}) (err error) {
	e.properties[key] = property
	return
}

func (e *SimpleEdge) Properties() (properties map[string]interface{}, err error) {
	properties = e.properties
	return
}

type SimpleGraph struct {
	vertices   map[Key]*SimpleVertex
	edges      map[Key]*SimpleEdge
	properties map[string]interface{}
}

func NewSimpleGraph() (g *SimpleGraph) {
	g = &SimpleGraph{
		vertices:   map[Key]*SimpleVertex{},
		edges:      map[Key]*SimpleEdge{},
		properties: map[string]interface{}{},
	}
	return
}

func (g *SimpleGraph) CreateVertex(key Key) (v Vertex, err error) {
	v = &SimpleVertex{
		key:        key,
		properties: map[string]interface{}{},
	}
	g.vertices[key] = v.(*SimpleVertex)
	return

}

func (g *SimpleGraph) CreateEdge(key Key, label string, head, tail Key) (e Edge, err error) {
	hv, ok := g.vertices[head]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	tv, ok := g.vertices[tail]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	se := &SimpleEdge{
		key:        key,
		label:      label,
		head:       hv,
		tail:       tv,
		properties: map[string]interface{}{},
	}
	g.edges[key] = se
	hv.incoming = append(hv.incoming, se)
	tv.outgoing = append(tv.outgoing, se)

	e = se
	return
}

func (g *SimpleGraph) RemoveVertex(key Key) (err error) {
	v, ok := g.vertices[key]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	for _, e := range v.incoming {
		err = g.RemoveEdge(e.key)
		if err != nil {
			return
		}
	}
	for _, e := range v.outgoing {
		err = g.RemoveEdge(e.key)
		if err != nil {
			return
		}
	}
	delete(g.vertices, key)
	return
}

func (g *SimpleGraph) RemoveEdge(key Key) (err error) {
	e, ok := g.edges[key]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	h := e.head
	tmp := h.incoming
	h.incoming = nil
	for _, i := range tmp {
		if i.key != key {
			h.incoming = append(h.incoming, i)
		}
	}
	t := e.tail
	tmp = t.outgoing
	t.outgoing = nil
	for _, o := range tmp {
		if o.key != key {
			t.outgoing = append(t.outgoing, o)
		}
	}
	delete(g.edges, key)
	return
}

func (g *SimpleGraph) GetVertex(key Key) (v Vertex, err error) {
	v, ok := g.vertices[key]
	if !ok {
		err = ErrKeyNotFound
	}
	return
}

func (g *SimpleGraph) GetEdge(key Key) (e Edge, err error) {
	e, ok := g.edges[key]
	if !ok {
		err = ErrKeyNotFound
	}
	return
}

func (g *SimpleGraph) WalkVertices(w VertexWalker) (err error) {
	for _, v := range g.vertices {
		w(v)
	}
	return
}

func (g *SimpleGraph) WalkEdges(w EdgeWalker) (err error) {
	for _, v := range g.edges {
		w(v)
	}
	return
}

func (g *SimpleGraph) GetProperty(key string) (property interface{}, err error) {
	property, ok := g.properties[key]
	if !ok {
		err = ErrKeyNotFound
		return
	}
	return
}

func (g *SimpleGraph) SetProperty(key string, property interface{}) (err error) {
	g.properties[key] = property
	return
}

func (g *SimpleGraph) Properties() (properties map[string]interface{}, err error) {
	properties = g.properties
	return
}
