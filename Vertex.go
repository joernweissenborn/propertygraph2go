package propertygraph2go

type Vertex struct {
	Id         string
	Incoming   []*Edge
	Outgoing   []*Edge
	Properties interface{}
}

func (v *Vertex) RemoveOutgoingEdge(id string) {
	for i, e := range v.Outgoing {
		if e.Id == id {
			v.Outgoing[i] = v.Outgoing[len(v.Outgoing)-1]
			v.Outgoing = v.Outgoing[:len(v.Outgoing)-1]
		}
	}
}

func (v *Vertex) RemoveIncomingEdge(id string) {
	for i, e := range v.Incoming {
		if e.Id == id {
			v.Incoming[i] = v.Incoming[len(v.Incoming)-1]
			v.Incoming = v.Incoming[:len(v.Incoming)-1]
		}
	}
}
