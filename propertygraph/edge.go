package propertygraph

type Edge interface {
	GraphItem
	Label()      string
	Head()       Vertex
	Tail()       Vertex
}
