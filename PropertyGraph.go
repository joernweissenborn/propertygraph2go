package PropertyGraph2Go

type PropertyGraph struct {
	vertices map[string]*Vertex
	edges    map[string]*Edge
}

type Vertex struct {
	Id         string
	Incoming   []*Edge
	Outgoing   []*Edge
	Properties interface{}
}

type Edge struct {
	Id         string
	Label      string
	Head       *Vertex
	Tail       *Vertex
	Properties interface{}
}

func New() *PropertyGraph {
	vs := make(map[string]*Vertex)
	es := make(map[string]*Edge)
	return &PropertyGraph{
		vs,
		es}
}

func (v *Vertex) RemoveOutgoingEdge(id string) {
	for i, e := range v.Outgoing {
		if e.Id == id {
			v.Outgoing[i] = v.Outgoing[len(v.Outgoing)-1]
			v.Outgoing = v.Outgoing[:len(v.Outgoing)-1]
		}
	}
}

func (v *Vertex) RemoveIncomingEdge(id string) {
	for i, e := range v.Incoming {
		if e.Id == id {
			v.Incoming[i] = v.Incoming[len(v.Incoming)-1]
			v.Incoming = v.Incoming[:len(v.Incoming)-1]
		}
	}
}

func (pg *PropertyGraph) CreateVertex(id string, properties interface{}) *Vertex {

	var nv Vertex

	nv.Id = id

	nv.Properties = properties

	pg.vertices[id] = &nv

	return &nv

}

func (pg *PropertyGraph) CreateEdge(id string, label string, head *Vertex, tail *Vertex, properties interface{}) *Edge {

	ne := Edge{
		id,
		label,
		head,
		tail,
		properties}
	pg.edges[id] = &ne
	head.Incoming = append(head.Incoming, &ne)
	tail.Outgoing = append(tail.Outgoing, &ne)

	return &ne

}

func (pg *PropertyGraph) RemoveVertex(id string) {
	if v := pg.GetVertex(id); v != nil {
		for _, e := range v.Incoming {
			pg.RemoveEdge(e.Id)
		}
		for _, e := range v.Outgoing {
			pg.RemoveEdge(e.Id)
		}
		delete(pg.vertices, id)
	}
}

func (pg *PropertyGraph) RemoveEdge(id string) {
	if e := pg.GetEdge(id); e != nil {
		e.Head.RemoveIncomingEdge(id)
		e.Tail.RemoveOutgoingEdge(id)
	}
	delete(pg.edges, id)
}

func (pg *PropertyGraph) GetVertex(id string) *Vertex {
	return pg.vertices[id]
}

func (pg *PropertyGraph) GetEdge(id string) *Edge {
	return pg.edges[id]
}

func (pg *PropertyGraph) GetIncomingEdgesByLabel(id string, label string) []*Edge {
	var es []*Edge // == nil
	for _, e := range pg.vertices[id].Incoming {
		if e.Label == label {
			es = append(es, e)
		}
	}
	return es

}

func (pg *PropertyGraph) GetOutgoingEdgesByLabel(id string, label string) []*Edge {
	var es []*Edge // == nil
	for _, e := range pg.vertices[id].Outgoing {
		if e.Label == label {
			es = append(es, e)
		}
	}
	return es
}
