package propertygraph

type Vertex interface {
	GraphItem
	Incoming() []Edge
	Outgoing()[]Edge
	RemoveIncomingEdge(id string)
	RemoveOutgoingEdge(id string)
	AddIncomingEdge(e Edge)
	AddOutgoingEdge(e Edge)
}

