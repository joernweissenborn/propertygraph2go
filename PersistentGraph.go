package propertygraph2go

type PersitentGraph struct {
	path string
	img *InMemoryGraph
	odg OnDiscGraph
}

func (pg PersitentGraph) SetPath(path string) {
	pg.path = path
}

func (pg PersitentGraph) Init() {
	pg.odg.init(pg.path)
	pg.img,_ = pg.odg.CreateInMemoryGraph()
}

func (pg PersitentGraph) CreateVertex(id string, properties interface{}) *Vertex{
	nv := pg.img.CreateVertex(id,properties)
	pg.odg.WriteVertex(nv)
	return nv
}

func (pg PersitentGraph) RemoveVertex(id string){
	pg.odg.RemoveVertex(id)
	updatedVertices := []string{}
	for _, out := range pg.img.GetVertex(id).Outgoing {
		pg.odg.RemoveEdge(out.Id)
		updatedVertices = append(updatedVertices,out.Head.Id)
	}
	for _, in := range pg.img.GetVertex(id).Incoming {
		pg.odg.RemoveEdge(in.Id)
		updatedVertices = append(updatedVertices,in.Tail.Id)
	}
	pg.img.RemoveVertex(id)

	for _, uv := range updatedVertices {
		pg.odg.WriteVertex(pg.img.GetVertex(uv))
	}
}

func (pg PersitentGraph) GetVertex(id string) *Vertex {
	return pg.img.GetVertex(id)
}
func (pg PersitentGraph) CreateEdge(id string, label string, head *Vertex, tail *Vertex, properties interface{}) *Edge{
	ne := pg.img.CreateEdge(id,label,head,tail,properties)
	pg.odg.WriteEdge(ne)
	pg.odg.WriteVertex(head)
	pg.odg.WriteVertex(tail)
	return ne
}
func (pg PersitentGraph)RemoveEdge(id string){
	de := pg.img.GetEdge(id)
	head := de.Head
	tail := de.Tail
	pg.img.RemoveEdge(id)
	pg.odg.RemoveEdge(id)
	pg.odg.WriteVertex(head)
	pg.odg.WriteVertex(tail)
}
func (pg PersitentGraph)GetEdge(id string) *Edge {
	return pg.img.GetEdge(id)
}
func (pg PersitentGraph)GetIncomingEdgesByLabel(id string, label string) []*Edge{
	return pg.img.GetIncomingEdgesByLabel(id,label)
}
func (pg PersitentGraph)GetOutgoingEdgesByLabel(id string, label string) []*Edge{
	return pg.img.GetOutgoingEdgesByLabel(id,label)
}
