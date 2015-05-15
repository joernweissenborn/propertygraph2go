package transcator

import (
	"github.com/joernweissenborn/propertygraph2go/propertygraph"
	"github.com/joernweissenborn/propertygraph2go/inmemorygraph"
	"testing"
)

var TA Transcator = New(inmemorygraph.New())

func TestReadEmpty(t *testing.T) {
	transaction:= func(graph *Transaction) {
		ev := graph.GetVertex("empty")
		if ev == nil {
			t.Log("Empty vertex query returned", ev)
		} else {
			t.Error("Empty vertex query did not return nil")
		}
		ee := graph.GetEdge("empty")
		if ee == nil {
			t.Log("Empty edge query returned", ee)
		} else {
			t.Error("Empty edge query did not return nil, returned", ee == nil)
		}
	}
	TA.Transaction(transaction)
}

func TestWriteReadVertex(t *testing.T) {

	type testprops struct {
		Name string
	}
	var props testprops
	props.Name = "Habib"
	transaction:= func(graph *Transaction) {

		graph.CreateVertex("test", props)
	}
	TA.Transaction(transaction)

	transaction= func(graph *Transaction) {

		v := graph.GetVertex("test")
		if v==nil {
			t.Fatal("could not retrieve vertex")
			return
		}
		rprops, found := v.Properties().(testprops)
	if !found {
		t.Error("Write Read test type assertion failed")
		return
	}

	if rprops.Name != "Habib" {
		t.Error("Write Read test failed. Expected 'Habib' but got", rprops.Name)
		return
	}
	}
	TA.Transaction(transaction)

	t.Log("Write Read test passed...")

}


func TestWriteReadEdge(t *testing.T) {
	var v1, v2 propertygraph.Vertex
	transaction:= func(graph *Transaction) {


		v1 = graph.CreateVertex("1", nil)
		v2 = graph.CreateVertex("2", nil)

		graph.CreateEdge("1", "testlabel", v1, v2, nil)
	}
	TA.Transaction(transaction)

	transaction= func(graph *Transaction) {
		e := graph.GetEdge("1")
		v1 = graph.GetVertex("1")
		v2 = graph.GetVertex("2")
		if e == nil {
			t.Error("Could not retrieve edge")
			return
		}

		if e.Head().Id() != v1.Id() {
			t.Error("Could not retrieve head")
			return
		}

		if e.Tail().Id() != v2.Id() {
			t.Error("Could not retrieve tail")
			return
		}

		if v1.Incoming()[0].Tail().Id() != v2.Id() {
			t.Error("Could not retrieve tail")
			return
		}

		if v2.Outgoing()[0].Head().Id() != v1.Id() {
			t.Error("Could not retrieve head")
			return
		}

	}
	TA.Transaction(transaction)

}

