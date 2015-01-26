package propertygraph2go

import (
	"os"
	"encoding/gob"
	"path"
)

type WriteableEdge struct {
	path string
	Eid string
	Elabel string
	Ehead string
	Etail string
	Eproperties interface {}
}
func (e WriteableEdge) Id() string {
	return e.Eid
}

func (e *WriteableEdge) readFromFile(Eid string) bool {
	file,err := os.OpenFile(path.Join(e.path,"edges", Eid),os.O_RDWR,0666)

	if err != nil {
		return false
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(e)
	return true
}

func (e WriteableEdge) Label() string {
	return e.Elabel

}
func (e WriteableEdge) Head() Vertex {
	var wv WriteableVertex
	wv.path = e.path
	wv.readFromFile(e.Ehead)
	return &wv
}
func (e WriteableEdge) Tail() Vertex {
	var wv WriteableVertex
	wv.path = e.path
	wv.readFromFile(e.Etail)
	return &wv

}
func (e WriteableEdge) Properties() interface {} {
	return e.Eproperties

}
