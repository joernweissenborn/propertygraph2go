package propertygraph2go

type SemiPersistentGraph struct {
	path string
	nonpers *InMemoryGraph
	pers OnDiscGraph
}

func (spg *SemiPersistentGraph) SetPath(path string) {
	spg.path = path
}

func (spg *SemiPersistentGraph) Init() (err error){
	spg.pers = OnDiscGraph{}
	spg.pers.init(spg.path)
	spg.nonpers,err = spg.pers.CreateInMemoryGraph()
	return
}

func (spg *SemiPersistentGraph) CreateVertex(id string, properties interface{}) *Vertex{
	nv := spg.nonpers.CreateVertex(id,properties)
	return nv
}
func (spg *SemiPersistentGraph) CreatePersistentVertex(id string, properties interface{}) *Vertex{
	nv := spg.nonpers.CreateVertex(id,properties)
	spg.pers.CreateVertex(id,properties)

	return nv
}

func (spg *SemiPersistentGraph) RemoveVertex(id string){
	spg.pers.RemoveVertex(id)
	spg.nonpers.RemoveVertex(id)
}

func (spg *SemiPersistentGraph) GetVertex(id string) *Vertex {
	return spg.nonpers.GetVertex(id)
}
func (spg *SemiPersistentGraph) CreateEdge(id string, label string, head *Vertex,
	tail *Vertex, properties interface{}) *Edge{
	_,err := spg.pers.GetVertex(head.Id)
	if err != nil {
		return spg.nonpers.CreateEdge(id,label,head,tail,properties)
	}
	_,err = spg.pers.GetVertex(tail.Id)
	if err != nil {
		return spg.nonpers.CreateEdge(id,label,head,tail,properties)
	}
	spg.pers.CreateEdge(id,label,head.Id,tail.Id,properties)
	return spg.nonpers.CreateEdge(id,label,head,tail,properties)
}

func (spg *SemiPersistentGraph)PersistVertex(id string){
	if !spg.isVertexPersistent(id) {
		v := spg.nonpers.GetVertex(id)
		if v == nil {
			return
		}
		spg.pers.CreateVertex(id,v.Properties)
		for _, e := range v.Incoming{
			if spg.isVertexPersistent(e.Tail.Id){
				spg.pers.CreateEdge(e.Id,e.Label,e.Head.Id,e.Tail.Id,e.Properties)
			}
		}
		for _, e := range v.Outgoing{
			if spg.isVertexPersistent(e.Head.Id){
				spg.pers.CreateEdge(e.Id,e.Label,e.Head.Id,e.Tail.Id,e.Properties)
			}
		}
	}
}

func (spg *SemiPersistentGraph)isVertexPersistent(id string) bool {
	_, err := spg.pers.GetVertex(id)
	if err != nil {
		return true
	}
	return false

}
func (spg *SemiPersistentGraph)RemoveEdge(id string){
	spg.pers.RemoveEdge(id)
	spg.nonpers.RemoveEdge(id)
}
func (spg *SemiPersistentGraph)GetEdge(id string) *Edge {
	return spg.nonpers.GetEdge(id)
}
func (spg *SemiPersistentGraph)GetIncomingEdgesByLabel(id string, label string) []*Edge{
	return spg.nonpers.GetIncomingEdgesByLabel(id,label)
}
func (spg *SemiPersistentGraph)GetOutgoingEdgesByLabel(id string, label string) []*Edge{
	return spg.nonpers.GetOutgoingEdgesByLabel(id,label)
}
