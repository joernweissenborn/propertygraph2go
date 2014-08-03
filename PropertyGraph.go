package propertygraph2go



type InMemoryGraph struct {
	vertices map[string]*Vertex
	edges    map[string]*Edge
}






//NewPersistant initializes a persistent Property graph at given location

func (pg *InMemoryGraph) CreateVertex(id string, properties interface{}) *Vertex {

	var nv Vertex

	nv.Id = id

	nv.Properties = properties

	pg.vertices[id] = &nv

	return &nv

}

func (pg *InMemoryGraph) CreateEdge(id string, label string, head *Vertex, tail *Vertex, properties interface{}) *Edge {

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

func (pg *InMemoryGraph) RemoveVertex(id string) {
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

func (pg *InMemoryGraph) RemoveEdge(id string) {
	if e := pg.GetEdge(id); e != nil {
		e.Head.RemoveIncomingEdge(id)
		e.Tail.RemoveOutgoingEdge(id)
	}
	delete(pg.edges, id)
}

func (pg *InMemoryGraph) GetVertex(id string) *Vertex {
	return pg.vertices[id]
}

func (pg *InMemoryGraph) GetEdge(id string) *Edge {
	return pg.edges[id]
}

func (pg *InMemoryGraph) GetIncomingEdgesByLabel(id string, label string) []*Edge {
	var es []*Edge // == nil
	for _, e := range pg.vertices[id].Incoming {
		if e.Label == label {
			es = append(es, e)
		}
	}
	return es

}

func (pg *InMemoryGraph) GetOutgoingEdgesByLabel(id string, label string) []*Edge {
	var es []*Edge // == nil
	for _, e := range pg.vertices[id].Outgoing {
		if e.Label == label {
			es = append(es, e)
		}
	}
	return es
}
