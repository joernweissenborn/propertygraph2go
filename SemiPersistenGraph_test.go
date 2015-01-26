package propertygraph2go

import (
	"os"
	"testing"
)

var semigraph SemiPersistentGraph
var semitestpath = "c:/semitest"

func TestSemiGraphInit(T *testing.T) {
	if err :=  os.RemoveAll(semitestpath); err != nil {
		T.Error(err)
	}
	semigraph.SetPath(semitestpath)
	 semigraph.Init()
}

func TestSemiGraphReInit(T *testing.T) {
	v1:=semigraph.CreatePersistentVertex("1",42)
	v2:=semigraph.CreatePersistentVertex("2","mario")
	v3:=semigraph.CreateVertex("3",nil)
	semigraph.CreateEdge("1","test",v1,v2,true)
	semigraph.CreateEdge("2","test",v1,v3,true)

	semigraph = SemiPersistentGraph{}
	semigraph.SetPath(semitestpath)
	semigraph.Init()

	v1 = semigraph.GetVertex("1")
	if v1 == nil {
		T.Error("Could not retrieve vertex")
	}
	v2 = semigraph.GetVertex("2")
	if v2 == nil {
		T.Error("Could not retrieve vertex")
	}
	v3 = semigraph.GetVertex("3")
	if v2 == nil {
		T.Error("Could retrieve vertex")
	}
	e := semigraph.GetEdge("1")
	if e == nil {
		T.Error("Could not retrieve edge")
	}
	e = semigraph.GetEdge("2")
	if e != nil {
		T.Error("Could retrieve vertex")
	}
}

