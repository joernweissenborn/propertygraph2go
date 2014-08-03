package propertygraph2go

type PropertyGraph interface {
	Init() error
	CreateVertex(id string, properties interface{}) *Vertex
	RemoveVertex(id string)
	GetVertex(id string) *Vertex
	CreateEdge(id string, label string, head *Vertex, tail *Vertex, properties interface{}) *Edge
	RemoveEdge(id string)
	GetEdge(id string) *Edge
	GetIncomingEdgesByLabel(id string, label string) []*Edge
	GetOutgoingEdgesByLabel(id string, label string) []*Edge
}

//New initialize a new non persistent graph
func New() *InMemoryGraph {
	vs := make(map[string]*Vertex)
	es := make(map[string]*Edge)
	return &InMemoryGraph{
		vs,
		es}
}

//NewSemiPersistent initializes a new semi persistent graph
func NewSemiPersistent(Path string) (*SemiPersistentGraph,error) {
	var semigraph SemiPersistentGraph
	semigraph.SetPath(Path)
	err :=  semigraph.Init()
	return &semigraph, err
}
