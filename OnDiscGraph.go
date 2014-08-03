package propertygraph2go

import (
	"errors"
	"os"
	"encoding/gob"
	"path"
	"path/filepath"
)

type OnDiscGraph struct {
	vertexpath string
	edgepath string
}

type WriteableVertex struct {
	Id string
	Incoming []string
	Outgoing []string
	Properties interface {}
}

func (v *WriteableVertex) AddOutgoingEdge(edgeId string) {
	v.Outgoing = append(v.Outgoing,edgeId)
}

func (v *WriteableVertex) RemoveOutgoingEdge(edgeId string) {
	for i, e := range v.Outgoing {
		if e == edgeId {
			v.Outgoing[i] = v.Outgoing[len(v.Outgoing)-1]
			v.Outgoing = v.Outgoing[:len(v.Outgoing)-1]
		}
	}
}

func (v *WriteableVertex) AddIncomingEdge(edgeId string) {
	v.Incoming = append(v.Incoming,edgeId)
}


func (v *WriteableVertex) RemoveIncomingEdge(edgeId string) {
	for i, e := range v.Incoming {
		if e == edgeId {
			v.Incoming[i] = v.Incoming[len(v.Incoming)-1]
			v.Incoming = v.Incoming[:len(v.Incoming)-1]
		}
	}
}

type WriteableEdge struct {
	Id string
	Label string
	Head string
	Tail string
	Properties interface {}
}

func (odg *OnDiscGraph) init(basepath string) error{
	odg.vertexpath = path.Join(basepath,"vertices")
	odg.edgepath = path.Join(basepath,"edges")
	err := os.MkdirAll(odg.vertexpath,os.ModeDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(odg.edgepath,os.ModeDir)
	if err != nil {
		return err
	}

	gob.Register(WriteableVertex{})
	gob.Register(WriteableEdge{})

	return nil

}


func (odg OnDiscGraph) WriteVertex(v *Vertex) error{
	if v==nil {
		return errors.New("Vertex is nil")
	}
	vertexpath := path.Join(odg.vertexpath,v.Id)
	return writeElement(convertToWritableVertex(v),vertexpath)

}

func (odg OnDiscGraph) WriteEdge(e *Edge) error{
	if e==nil {
		return errors.New("Edge is nil")
	}
	edgepath := path.Join(odg.edgepath,e.Id)

	return writeElement(convertToWritableEdge(e),edgepath)

}

func writeElement(element interface {},path string) error{

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)

	err = enc.Encode(element)
	if err != nil {
		return err
	}

	return nil
}
func (odg OnDiscGraph) CreateVertex(id string, properties interface{}) (v *WriteableVertex, err error){
	wv := WriteableVertex{
		id,
		[]string{},
		[]string{},
		properties}
	err = writeElement(wv,path.Join(odg.vertexpath,id))
	v = &wv
	return
}

func (odg OnDiscGraph) RemoveVertex(id string) (err error) {
	v, err := odg.GetVertex(id)
	if err!= nil {
		return err
	}
	for _,e := range v.Incoming {
		err := odg.RemoveEdge(e)
		if err!= nil {
			return err
		}
	}
	return os.Remove(path.Join(odg.vertexpath,id))
}

func (odg OnDiscGraph) CreateEdge(id string, label string,
	head string, tail string, properties interface{})(e WriteableEdge, err error){
	e = WriteableEdge{
		id,
		label,
		head,
		tail,
		properties}
	v,err := odg.GetVertex(head)
	if err != nil {
		return
	}
	v.AddIncomingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id))
	if err != nil {
		return
	}
	v,err = odg.GetVertex(tail)
	if err != nil {
		return
	}
	v.AddOutgoingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id))
	if err != nil {
		return
	}
	err = writeElement(e,path.Join(odg.edgepath,id))

	return
}


func (odg OnDiscGraph) RemoveEdge(id string) (err error) {
	e, err := odg.GetEdge(id)
	if err!= nil {
		return err
	}
	v, err := odg.GetVertex(e.Head)
	if err!= nil {
		return err
	}
	v.RemoveIncomingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id))
	if err!= nil {
		return err
	}
	v, err = odg.GetVertex(e.Tail)
	if err!= nil {
		return err
	}
	v.RemoveOutgoingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id))
	if err!= nil {
		return err
	}
	return os.Remove(path.Join(odg.edgepath,id))
}


func (odg OnDiscGraph) GetVertex(id string) (v WriteableVertex, err error) {
	vertexpath := path.Join(odg.vertexpath,id)
	file,err := os.Open(vertexpath)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&v)
	return
}

func (odg OnDiscGraph) GetEdge(id string) (e WriteableEdge, err error) {
	edgepath := path.Join(odg.edgepath,id)
	file,err := os.Open(edgepath)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&e)
	return
}

func getItem(path string,item interface {})error{
	file,err := os.Open(path)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	return dec.Decode(&item)
}

func (odg OnDiscGraph) CreateInMemoryGraph() (img *InMemoryGraph, err error) {
	img = New()
	err = filepath.Walk(odg.vertexpath,func(p string, f os.FileInfo, err error) error{
		if !f.IsDir() {
			v, err := odg.GetVertex(f.Name())
			if err!=nil {
				return err
			}

			odg.addVertexToImg(img,v)

		}
		return nil
})
	return
}

func (odg OnDiscGraph) addVertexToImg(img *InMemoryGraph, v WriteableVertex) {
	imv := img.GetVertex(v.Id)
	if imv == nil{
		imv = img.CreateVertex(v.Id,v.Properties)
	}
	for _,in := range v.Incoming {
		e, _ := odg.GetEdge(in)

		odg.addEdgeToImg(img,e)
	}
}

func (odg OnDiscGraph) addEdgeToImg(img *InMemoryGraph, e WriteableEdge) {

	ime := img.GetEdge(e.Id)
	if ime == nil {
		tail := img.GetVertex(e.Tail)
		if tail == nil {
			v,_ := odg.GetVertex(e.Tail)
			odg.addVertexToImg(img,v)
			tail = img.GetVertex(e.Tail)
		}
		head := img.GetVertex(e.Head)
		if head == nil {
			v,_ := odg.GetVertex(e.Head)
			odg.addVertexToImg(img,v)
			head = img.GetVertex(e.Head)
		}
		img.CreateEdge(e.Id,e.Label,head,tail,e.Properties)
	}
}

func flushFile(f *os.File) *os.File{
	path := f.Name()
	f.Close()
	os.Remove(path)
	f,_ = os.Create(path)
	return f
}


func convertToWritableVertex(v *Vertex) (wv WriteableVertex) {
	wv.Id = v.Id
	for _,e := range v.Incoming {
		wv.Incoming = append(wv.Incoming,e.Id)
	}
	for _,e := range v.Outgoing {
		wv.Outgoing = append(wv.Outgoing, e.Id)
	}
	wv.Properties = v.Properties
	return
}
func convertToWritableEdge(e *Edge) (we WriteableEdge) {
	we.Id = e.Id
	we.Label = e.Label
	we.Head = e.Tail.Id
	we.Tail = e.Tail.Id
	we.Properties = e.Properties
	return
}
