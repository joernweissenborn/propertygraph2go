package propertygraph2go
/*
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
	id string
	incoming []*WriteableEdge
	outgoing []*WriteableEdge
	properties interface {}
}

func (v WriteableVertex) Id() string {
	return v.id
	
}

func (v WriteableVertex) Outgoing() ([]Edge) {

	return v.outgoing
	
}
func (v WriteableVertex) Incoming() ([]Edge) {

	return v.incoming

}
func (v WriteableVertex) Properties() (interface {}) {

	return v.properties

}

func (v *WriteableVertex) AddOutgoingEdge(e *WriteableEdge) {
	v.outgoing = append(v.outgoing,e)
}

func (v *WriteableVertex) RemoveOutgoingEdge(edgeId string) {
	for i, e := range v.outgoing {
		if e.Id() == edgeId {
			v.outgoing[i] = v.outgoing[len(v.outgoing)-1]
			v.outgoing = v.outgoing[:len(v.outgoing)-1]
		}
	}
}

func (v *WriteableVertex) AddIncomingEdge(e *WriteableEdge) {
	v.incoming = append(v.incoming,e)
}


func (v *WriteableVertex) RemoveIncomingEdge(edgeId string) {
	for i, e := range v.incoming {
		if e.Id() == edgeId {
			v.incoming[i] = v.incoming[len(v.incoming)-1]
			v.incoming = v.incoming[:len(v.incoming)-1]
		}
	}
}

type WriteableEdge struct {
	id string
	label string
	head *WriteableVertex
	tail *WriteableVertex
	properties interface {}
}
func (e WriteableEdge) Id() string {
	return e.id

}
func (e WriteableEdge) Label() string {
	return e.label

}
func (e WriteableEdge) Head() Vertex {
	return e.head

}
func (e WriteableEdge) Tail() Vertex {
	return e.tail

}
func (e WriteableEdge) Properties() interface {} {
	return e.properties

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


func (odg OnDiscGraph) WriteVertex(v Vertex) error{
	if v==nil {
		return errors.New("Vertex is nil")
	}
	vertexpath := path.Join(odg.vertexpath,v.Id())
	return writeElement(convertToWritableVertex(v),vertexpath)

}

func (odg OnDiscGraph) WriteEdge(e Edge) error{
	if e==nil {
		return errors.New("Edge is nil")
	}
	edgepath := path.Join(odg.edgepath,e.Id())

	return writeElement(convertToWritableEdge(e),edgepath)

}

func writeElement(element interface {},path string) error{
	os.Remove(path)
	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666)
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
		[]*WriteableEdge{},
		[]*WriteableEdge{},
		properties}
	err = writeElement(wv,path.Join(odg.vertexpath,id))
	v = &wv
	return
}

func (odg OnDiscGraph) RemoveVertex(id string) {
	v := odg.GetVertex(id)

	for _,e := range v.Incoming() {
		err := odg.RemoveEdge(e.Id())
		if err!= nil {
			return
		}
	}
	os.Remove(path.Join(odg.vertexpath,id))
}

func (odg OnDiscGraph) CreateEdge(id string, label string,
	head Vertex, tail Vertex, properties interface{})(e Edge){
	te := WriteableEdge{
		id,
		label,
		convertToWritableVertex(head),
		convertToWritableVertex(tail),
		properties}
	v := convertToWritableVertex(head)
	v.AddIncomingEdge(&te)
	e = &te
	err := writeElement(v,path.Join(odg.vertexpath,v.Id()))
	if err != nil {
		return
	}
	v = convertToWritableVertex(tail)

	v.AddOutgoingEdge(&te)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id()))
	if err != nil {
		return
	}
	err = writeElement(e,path.Join(odg.edgepath,id))

	return
}


func (odg OnDiscGraph) RemoveEdge(id string) (err error) {
	e := odg.GetEdge(id)

	v := odg.GetVertex(e.Head().Id())

	v.RemoveIncomingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id()))
	if err!= nil {
		return err
	}
	v = odg.GetVertex(e.Tail().Id())

	v.RemoveOutgoingEdge(id)
	err = writeElement(v,path.Join(odg.vertexpath,v.Id()))
	if err!= nil {
		return err
	}
	return os.Remove(path.Join(odg.edgepath,id))
}


func (odg OnDiscGraph) GetVertex(id string) (v Vertex) {
	vertexpath := path.Join(odg.vertexpath,id)
	file,err := os.OpenFile(vertexpath,os.O_RDWR,0666)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&v)
	return
}

func (odg OnDiscGraph) GetEdge(id string) (e Edge) {
	edgepath := path.Join(odg.edgepath,id)
	file,err := os.OpenFile(edgepath,os.O_RDWR,0666)
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
	img = NewInMemoryGraph()
	err = filepath.Walk(odg.vertexpath,func(p string, f os.FileInfo, err error) error{
		if !f.IsDir() {
			v := odg.GetVertex(f.Name())
			
			odg.addVertexToImg(img,v)

		}
		return nil
})
	return
}

func (odg OnDiscGraph) addVertexToImg(img *InMemoryGraph, v Vertex) {
	imv := img.GetVertex(v.Id())
	if imv == nil{
		imv = img.CreateVertex(v.Id(),v.Properties)
	}
	for _,in := range v.Incoming() {
		e := odg.GetEdge(in.Id())

		odg.addEdgeToImg(img,e)
	}
}

func (odg OnDiscGraph) addEdgeToImg(img *InMemoryGraph, e Edge) {

	ime := img.GetEdge(e.Id())
	if ime == nil {
		tail := img.GetVertex(e.Tail().Id())
		if tail == nil {
			v := odg.GetVertex(e.Tail().Id())
			odg.addVertexToImg(img,v)
			tail = img.GetVertex(e.Tail().Id())
		}
		head := img.GetVertex(e.Head().Id())
		if head == nil {
			v := odg.GetVertex(e.Head().Id())
			odg.addVertexToImg(img,v)
			head = img.GetVertex(e.Head().Id())
		}
		img.CreateEdge(e.Id(),e.Label(),head,tail,e.Properties())
	}
}

func flushFile(f *os.File) *os.File{
	path := f.Name()
	f.Close()
	os.Remove(path)
	f,_ = os.Create(path)
	return f
}


func convertToWritableVertex(v Vertex) (wv *WriteableVertex) {
	var t WriteableVertex
	
	wv = &t 
	wv.id = v.Id()
	for _,e := range v.Incoming() {
		wv.incoming = append(wv.incoming,e)
	}
	for _,e := range v.Outgoing() {
		wv.outgoing = append(wv.outgoing, e)
	}
	wv.properties = v.Properties()
	return
}
func convertToWritableEdge(e Edge) (we *WriteableEdge) {
	var t WriteableEdge
	we = &t
	we.id = e.Id()
	we.label = e.Label()
	we.head = e.Tail().Id()
	we.tail = e.Tail().Id()
	we.properties = e.Properties()
	return
}
*/
