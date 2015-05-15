package inmemorygraph

import "github.com/joernweissenborn/propertygraph2go/propertygraph"

//New initialize a new non persistent graph
func New() *InMemoryGraph {
	var g InMemoryGraph
	g.Init()
	return &g
}

type InMemVertex struct {
	id         string
	incoming   []*InMemEdge
	outgoing   []*InMemEdge
	properties interface{}
}

func (v *InMemVertex) Id() string {
	return v.id
}

func (v *InMemVertex) Incoming() []propertygraph.Edge {
	es := []propertygraph.Edge{}
	for _, e:= range v.incoming {
		es = append(es,e)
		
	}
	return es
}

func (v *InMemVertex) Outgoing() []propertygraph.Edge {
	es := []propertygraph.Edge{}
	for _, e:= range v.outgoing {
		es = append(es,e)

	}
	return es}

func (v *InMemVertex) Properties() interface {} {
	return v.properties
}

func (v *InMemVertex) RemoveOutgoingEdge(id string) {
	for i, e := range v.outgoing {
		if e.Id() == id {
			v.outgoing[i] = v.outgoing[len(v.outgoing)-1]
			v.outgoing = v.outgoing[:len(v.outgoing)-1]
		}
	}
}

func (v *InMemVertex) RemoveIncomingEdge(id string) {
	for i, e := range v.incoming {
		if e.Id() == id {
			v.incoming[i] = v.incoming[len(v.incoming)-1]
			v.incoming = v.incoming[:len(v.incoming)-1]
		}
	}
}

func (v *InMemVertex) AddIncomingEdge(e propertygraph.Edge) {
	 v.incoming = append(v.incoming,e.(*InMemEdge))
}


func (v *InMemVertex) AddOutgoingEdge(e propertygraph.Edge) {
	 v.outgoing = append(v.outgoing,e.(*InMemEdge))
}



type InMemEdge struct {
	id         string
	label      string
	head       *InMemVertex
	tail       *InMemVertex
	properties interface{}
}

func (e *InMemEdge) Id() string {
	return e.id
}

func (e *InMemEdge) Label() string {
	return e.label
}

func (e *InMemEdge) Head() propertygraph.Vertex {
	return e.head
}

func (e *InMemEdge) Tail() propertygraph.Vertex {
	return e.tail
}

func (e *InMemEdge) Properties() interface {} {
	return e.properties
}




type InMemoryGraph struct {
	vertices map[string]*InMemVertex
	edges    map[string]*InMemEdge
}



func (pg *InMemoryGraph) Init() {


	pg.vertices = make(map[string]*InMemVertex)
	pg.edges = make(map[string]*InMemEdge)

}

func (pg *InMemoryGraph) CreateVertex(id string, properties interface{}) propertygraph.Vertex {

	var nv InMemVertex

	nv.id = id

	nv.properties = properties

	pg.vertices[id] = &nv

	return &nv

}

func (pg *InMemoryGraph) CreateEdge(id string, label string, head propertygraph.Vertex, tail propertygraph.Vertex, properties interface{}) propertygraph.Edge {
	hv := pg.GetVertex(head.Id()).(*InMemVertex)
	tv := pg.GetVertex(tail.Id()).(*InMemVertex)
	ne := InMemEdge{
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

func (pg *InMemoryGraph) GetVertex(id string) propertygraph.Vertex {
	if v := pg.vertices[id]; v !=nil {
		return v
	}
	return nil
}

func (pg *InMemoryGraph) GetEdge(id string) propertygraph.Edge {
	if v := pg.edges[id]; v != nil {
		return v
	}
	return nil
}

func (pg *InMemoryGraph) GetIncomingEdgesByLabel(id string, label string) []propertygraph.Edge {
	var es []propertygraph.Edge // == nil
	v := pg.vertices[id]
	if v != nil {
		for _, e := range v.Incoming() {
			if e.Label() == label {
				es = append(es, e)
			}
		}
	}
	return es

}

func (pg *InMemoryGraph) GetOutgoingEdgesByLabel(id string, label string) []propertygraph.Edge {
	var es []propertygraph.Edge // == nil
	v := pg.vertices[id]
	if v != nil {
		for _, e := range v.Outgoing() {
			if e.Label() == label {
				es = append(es, e)
			}
		}
	}
	return es
}
