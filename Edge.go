package propertygraph2go




type Edge interface {
	Id()         string
	Label()      string
	Head()       Vertex
	Tail()       Vertex
	Properties() interface{}
}
