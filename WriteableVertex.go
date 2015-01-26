package propertygraph2go

import (
	"encoding/gob"
	"os"
	"path"
)

type WriteableVertex struct {
	path string
	Vid string
	Vincoming []string
	Voutgoing []string
	Vproperties interface {}
}

func (v *WriteableVertex) readFromFile(id string) bool {
	file,err := os.OpenFile(path.Join(v.path,"vertices", id),os.O_RDWR,0666)
	if err != nil {
		return false
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&v)
	return true
}

func (v *WriteableVertex) Id() string {
	return v.Vid

}

func (v *WriteableVertex) Outgoing() ([]Edge) {
	edges := []Edge{}
	for _, id := range v.Voutgoing {
		var e WriteableEdge
		e.path = v.path
		e.readFromFile(id)
		edges = append(edges, &e)
	}
	return edges
}


func (v *WriteableVertex) Incoming() ([]Edge) {

	edges := []Edge{}
	for _, id := range v.Vincoming {

		var e WriteableEdge
		e.path = v.path
		e.readFromFile(id)
		edges = append(edges, &e)
	}
	return edges
}
func (v *WriteableVertex) Properties() (interface {}) {

	return v.Vproperties

}

func (v *WriteableVertex) AddOutgoingEdge(e Edge) {
	v.Voutgoing = append(v.Voutgoing,e.Id())
}

func (v *WriteableVertex) RemoveOutgoingEdge(edgeId string) {
	for i, id := range v.Voutgoing {
		if id == edgeId {
			v.Voutgoing[i] = v.Voutgoing[len(v.Voutgoing)-1]
			v.Voutgoing = v.Voutgoing[:len(v.Voutgoing)-1]
		}
	}
}

func (v *WriteableVertex) AddIncomingEdge(e Edge) {
	v.Vincoming = append(v.Vincoming,e.Id())
}


func (v *WriteableVertex) RemoveIncomingEdge(edgeId string) {
	for i, id := range v.Vincoming {
		if id == edgeId {
			v.Vincoming[i] = v.Vincoming[len(v.Vincoming)-1]
			v.Vincoming = v.Vincoming[:len(v.Vincoming)-1]
		}
	}
}
