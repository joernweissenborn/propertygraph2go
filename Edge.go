package propertygraph2go

type Edge struct {
	Id         string
	Label      string
	Head       *Vertex
	Tail       *Vertex
	Properties interface{}
}
