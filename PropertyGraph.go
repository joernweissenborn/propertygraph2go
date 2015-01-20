package propertygraph2go

//New initialize a new non persistent graph
func NewInMemoryGraph() *InMemoryGraph {
	var g InMemoryGraph
	g.Init()
	return &g
}

type inMemVertex struct {
	id         string
	incoming   []*inMemEdge
	outgoing   []*inMemEdge
	properties interface{}
}

func (v *inMemVertex) Id() string {
	return v.id
}

func (v *inMemVertex) Incoming() []Edge {
	es := []Edge{}
	for _, e:= range v.incoming {
		es = append(es,e)
		
	}
	return es
}

func (v *inMemVertex) Outgoing() []Edge {
	es := []Edge{}
	for _, e:= range v.outgoing {
		es = append(es,e)

	}
	return es}

func (v *inMemVertex) Properties() interface {} {
	return v.properties
}

func (v *inMemVertex) RemoveOutgoingEdge(id string) {
	for i, e := range v.outgoing {
		if e.Id() == id {
			v.outgoing[i] = v.outgoing[len(v.outgoing)-1]
			v.outgoing = v.outgoing[:len(v.outgoing)-1]
		}
	}
}

func (v *inMemVertex) RemoveIncomingEdge(id string) {
	for i, e := range v.incoming {
		if e.Id() == id {
			v.incoming[i] = v.incoming[len(v.incoming)-1]
			v.incoming = v.incoming[:len(v.incoming)-1]
		}
	}
}


type inMemEdge struct {
	id         string
	label      string
	head       *inMemVertex
	tail       *inMemVertex
	properties interface{}
}

func (e *inMemEdge) Id() string {
	return e.id
}

func (e *inMemEdge) Label() string {
	return e.label
}

func (e *inMemEdge) Head() Vertex {
	return e.head
}

func (e *inMemEdge) Tail() Vertex {
	return e.tail
}

func (e *inMemEdge) Properties() interface {} {
	return e.properties
}




type InMemoryGraph struct {
	vertices map[string]*inMemVertex
	edges    map[string]*inMemEdge
}



func (pg *InMemoryGraph) Init() {


	pg.vertices = make(map[string]*inMemVertex)
	pg.edges = make(map[string]*inMemEdge)

}

func (pg *InMemoryGraph) CreateVertex(id string, properties interface{}) Vertex {

	var nv inMemVertex

	nv.id = id

	nv.properties = properties

	pg.vertices[id] = &nv

	return &nv

}

func (pg *InMemoryGraph) CreateEdge(id string, label string, head Vertex, tail Vertex, properties interface{}) Edge {
	hv := pg.GetVertex(head.Id()).(*inMemVertex)
	tv := pg.GetVertex(tail.Id()).(*inMemVertex)
	ne := inMemEdge{
		id,
		label,
		hv,
		tv,
		properties}
	pg.edges[id] = &ne
	hv.incoming = append(hv.incoming, &ne)
	tv.outgoing = append(tv.outgoing, &ne)

	return &ne

}

func (pg *InMemoryGraph) RemoveVertex(id string) {
	if v := pg.GetVertex(id); v != nil {
		for _, e := range v.Incoming() {
			pg.RemoveEdge(e.Id())
		}
		for _, e := range v.Outgoing() {
			pg.RemoveEdge(e.Id())
		}
		delete(pg.vertices, id)
	}
}

func (pg *InMemoryGraph) RemoveEdge(id string) {
	if e := pg.GetEdge(id); e != nil {
		e.Head().RemoveIncomingEdge(id)
		e.Tail().RemoveOutgoingEdge(id)
	}
	delete(pg.edges, id)
}

func (pg *InMemoryGraph) GetVertex(id string) Vertex {
	if v := pg.vertices[id]; v !=nil {
		return v
	}
	return new(inMemVertex)
}

func (pg *InMemoryGraph) GetEdge(id string) Edge {
	return pg.edges[id]
}

func (pg *InMemoryGraph) GetIncomingEdgesByLabel(id string, label string) []Edge {
	var es []Edge // == nil
	for _, e := range pg.vertices[id].Incoming() {
		if e.Label() == label {
			es = append(es, e)
		}
	}
	return es

}

func (pg *InMemoryGraph) GetOutgoingEdgesByLabel(id string, label string) []Edge {
	var es []Edge // == nil
	for _, e := range pg.vertices[id].Outgoing() {
		if e.Label() == label {
			es = append(es, e)
		}
	}
	return es
}
