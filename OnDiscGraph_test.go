package propertygraph2go
/*
import (
	"testing"
	"os"
)

var testpath = "c:/test"

var odg OnDiscGraph

var graph = New()

var v1 = graph.CreateVertex("1",42)

var v2 = graph.CreateVertex("2",nil)

var e1 = graph.CreateEdge("1","test", v1,v2, nil)

func TestInit(T *testing.T){
	if err :=  os.RemoveAll(testpath); err != nil {
		T.Error(err)
	}
	if err :=  odg.init(testpath); err != nil {
		T.Error(err)
}
}

func TestWriteVertex(T *testing.T) {

	if err :=  odg.WriteVertex(v1); err != nil {
		T.Error(err)
	}
	if err :=  odg.WriteVertex(v2); err != nil {
		T.Error(err)
	}

}
func TestOverwriteVertex(T *testing.T) {

	if err :=  odg.WriteVertex(v1); err != nil {
		T.Error(err)
	}

}

func TestWriteEdge(T *testing.T) {

	if err :=  odg.WriteEdge(e1); err != nil {
		T.Error(err)
	}

}

func TestReadVertex(T *testing.T) {
	_, err :=  odg.GetVertex("non-existant")
	if  err == nil {
		T.Error("Could retrieve non-existant vertex")
	}
	v, err :=  odg.GetVertex("1")
	if  err != nil {
		T.Error(err)
	}
	if v.Properties != convertToWritableVertex(v1).Properties {
		T.Error("Could not retrieve vertex")
	}
}

func TestReadEdge(T *testing.T) {
	_, err :=  odg.GetEdge("non-existant")
	if  err == nil {
		T.Error("Could retrieve non-existant edge")
	}
	e, err :=  odg.GetEdge("1")
	if  err != nil {
		T.Error(err)
	}
	if e.Properties != convertToWritableEdge(e1).Properties {
		T.Error("Could not retrieve edge")
	}
}
func TestReadGraph(T *testing.T) {
	ng, err :=  odg.CreateInMemoryGraph()

	if  err != nil {
		T.Error(err)
	}

	v := ng.GetVertex(v1.Id)
	if  v == nil {
		T.Error("Could not retrieve vertex")
	}
	v = ng.GetVertex(v2.Id)
	if  v == nil {
		T.Error("Could not retrieve vertex")
	}
	e := ng.GetEdge(e1.Id)
	if  e == nil {
		T.Error("Could not retrieve edge")
	}
}



*/
