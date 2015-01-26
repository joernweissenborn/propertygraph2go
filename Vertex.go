package propertygraph2go
type Vertex interface {
	Id() string
	Incoming() []Edge
	Outgoing()[]Edge
	RemoveIncomingEdge(id string)
	RemoveOutgoingEdge(id string)
	Properties() interface{}
	AddIncomingEdge(e Edge)
	AddOutgoingEdge(e Edge)
}

