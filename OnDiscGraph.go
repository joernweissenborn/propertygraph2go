package propertygraph2go

import (
	"errors"
	"os"
	"encoding/gob"
	"path"
	"path/filepath"
	"fmt"
)

type OnDiscGraph struct {
	basepath string
	vertexpath string
	edgepath string
}



func (odg *OnDiscGraph) Init(basepath string) {
	odg.basepath = basepath
	odg.vertexpath = path.Join(basepath,"vertices")
	odg.edgepath = path.Join(basepath,"edges")
	err := os.MkdirAll(odg.vertexpath,os.ModeDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(odg.edgepath,os.ModeDir)
	if err != nil {
		panic(err)
	}

	gob.Register(WriteableVertex{})
	gob.Register(WriteableEdge{})

	return

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
func (odg OnDiscGraph) CreateVertex(id string, properties interface{}) Vertex{
	wv := WriteableVertex{
		odg.basepath,
		id,
		[]string{},
		[]string{},
		properties}
	
	writeElement(wv,path.Join(odg.vertexpath,id))

	return &wv
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
		odg.basepath,
		id,
		label,
		head.Id(),
		tail.Id(),
		properties}

	whead:= convertToWritableVertex(head)

	whead.AddIncomingEdge(&te)
	e = &te
	err := writeElement(whead,path.Join(odg.vertexpath,head.Id()))
	if err != nil {
		return
	}

		wtail := convertToWritableVertex(tail)
	wtail.AddOutgoingEdge(&te)
	err = writeElement(wtail,path.Join(odg.vertexpath,tail.Id()))
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


func (odg OnDiscGraph) GetVertex(id string) (Vertex) {
	var wv WriteableVertex
	wv.path = odg.basepath
	ok := wv.readFromFile(id)
	if !ok {
		return nil
		
	}
	return &wv
}

func (odg OnDiscGraph) GetEdge(id string) (Edge) {
	var we WriteableEdge
	we.path = odg.basepath
	ok := we.readFromFile(id)
	if !ok {
		return nil

	}

	return &we
}

func getItem(path string,item interface {})error{
	file,err := os.Open(path)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(file)
	return dec.Decode(&item)
}

func (odg *OnDiscGraph) CreateInMemoryGraph() (img *InMemoryGraph, err error) {
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
	fmt.Println("ggggggggggggggg",v)

	for _,e := range v.Incoming() {
		odg.addEdgeToImg(img,e)
	}	
	for _,e := range v.Outgoing() {
		
		odg.addEdgeToImg(img,e)
	}
}

func (odg OnDiscGraph) addEdgeToImg(img *InMemoryGraph, e Edge) {
	ime := img.GetEdge(e.Id())

	if ime == nil {
		tail := img.GetVertex(e.Tail().Id())
		if tail == nil {

			odg.addVertexToImg(img,e.Tail())
			tail = img.GetVertex(e.Tail().Id())
		}
		head := img.GetVertex(e.Head().Id())
		if head == nil {
			odg.addVertexToImg(img,e.Head())
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


func convertToWritableVertex(v Vertex) (*WriteableVertex) {
	var wv WriteableVertex
	wv.Vid = v.Id()
	for _,e := range v.Incoming() {
		wv.Vincoming = append(wv.Vincoming,e.Id())
	}
	for _,e := range v.Outgoing() {
		wv.Voutgoing = append(wv.Voutgoing, e.Id())
	}
	wv.Vproperties = v.Properties()
	return &wv
}
func convertToWritableEdge(e Edge) (*WriteableEdge) {
		var we WriteableEdge
	we.Eid = e.Id()
	we.Elabel = e.Label()
	we.Etail = e.Tail().Id()
	we.Ehead = e.Head().Id()
	we.Eproperties = e.Properties()
	return &we
}

