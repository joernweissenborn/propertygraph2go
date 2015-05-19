package propertygraph

type PropertyGraph interface {
	CreateVertex(id string, properties interface{}) Vertex
	RemoveVertex(id string)
	GetVertex(id string) Vertex
	CreateEdge(id string, label string, head Vertex, tail Vertex, properties interface{}) Edge
	RemoveEdge(id string)
	GetEdge(id string) Edge
}
