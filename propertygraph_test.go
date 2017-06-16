package propertygraph2go

import "testing"

func TestSimpleGraph(t *testing.T) {
	GraphTestSuite(t, NewSimpleGraph())
}

func GraphTestSuite(t *testing.T, g Graph) {
	testQueryNonExistent(t, g)
	testWriteReadVertex(t, g)
	testWriteReadEdge(t, g)
	testRemoveEdge(t, g)
	testRemoveVertex(t, g)
}

func testQueryNonExistent(t *testing.T, g Graph) {
	t.Log("Testing ")
	_, err := g.GetVertex("none")
	if err == nil {
		t.Fatal("Non existent vertex query returned no error")
	} else if err != ErrKeyNotFound {
		t.Fatal("Non existent vertex query did not return 'ErrKeyNotFound' but", err)
	}
	_, err = g.GetEdge("none")
	if err == nil {
		t.Fatal("Non existent edge query returned no error")
	} else if err != ErrKeyNotFound {
		t.Fatal("Non existent edge query did not return 'ErrKeyNotFound' but", err)
	}
}

func testWriteReadVertex(t *testing.T, g Graph) {
	v, err := g.CreateVertex("test")
	if err != nil {
		t.Fatal(err)
	}
	err = v.SetProperty("name", "dennis")
	if err != nil {
		t.Fatal(err)
	}

	v, err = g.GetVertex("test")
	if err != nil {
		t.Fatal(err)
	}

	name, err := v.GetProperty("name")
	if err != nil {
		t.Fatal(err)
	}

	if name.(string) != "dennis" {
		t.Fatal("Write Read test failed. Expected 'dennis' but got", name)
		return
	}

}

func testWriteReadEdge(t *testing.T, g Graph) {

	v1, err := g.CreateVertex("1")
	if err != nil {
		t.Fatal(err)
	}
	v2, err := g.CreateVertex("2")
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.CreateEdge("1", "testlabel", v1.Key(), v2.Key())
	if err != nil {
		t.Fatal(err)
	}

	e, err := g.GetEdge("1")
	if err != nil {
		t.Fatal(err)
	}

	if e.Head().Key() != v1.Key() {
		t.Fatal("Could not retrieve head")
	}

	if e.Tail().Key() != v2.Key() {
		t.Fatal("Could not retrieve tail")
	}

	if v1.Incoming()[0].Key() != e.Key() {
		t.Fatal("Could not retrieve edge")
		return
	}

	if v2.Outgoing()[0].Key() != e.Key() {
		t.Fatal("Could not retrieve edge")
	}

}

func testRemoveVertex(t *testing.T, g Graph) {

	v1, err := g.CreateVertex("1")
	if err != nil {
		t.Fatal(err)
	}
	v2, err := g.CreateVertex("2")
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.CreateEdge("1", "testlabel", v1.Key(), v2.Key())
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.CreateEdge("2", "testlabel", v2.Key(), v1.Key())
	if err != nil {
		t.Fatal(err)
	}

	err = g.RemoveVertex("1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.GetVertex("1")
	if err == nil {
		t.Fatal("Vertex was not removed")
	}

	_, err = g.GetEdge("1")
	if err == nil {
		t.Fatal("Incoming edge was not removed")
	}

	_, err = g.GetEdge("2")
	if err == nil {
		t.Fatal("Outgoing edge was not removed")
	}

	if len(v2.Incoming()) != 0 {
		t.Fatal("Incoming edge was not removed from vertex")
	}

	if len(v2.Outgoing()) != 0 {
		t.Fatal("Outgoing edge was not removed from vertex")
	}

}

func testRemoveEdge(t *testing.T, g Graph) {

	v1, err := g.CreateVertex("1")
	if err != nil {
		t.Fatal(err)
	}
	v2, err := g.CreateVertex("2")
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.CreateEdge("1", "testlabel", v1.Key(), v2.Key())
	if err != nil {
		t.Fatal(err)
	}

	err = g.RemoveEdge("1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = g.GetEdge("1")
	if err == nil {
		t.Fatal("Edge was not removed")
	}

	if len(v1.Incoming()) != 0 {
		t.Fatal("Incoming edge was not removed from vertex")
	}

	if len(v2.Outgoing()) != 0 {
		t.Fatal("Outgoing edge was not removed from vertex")
	}

}

func TestSerialize(t *testing.T) {
	g := NewSimpleGraph()

	v1, err := g.CreateVertex("1")
	if err != nil {
		t.Fatal(err)
	}
	sv1, err := EncodeVertexJSON(v1)
	if err != nil {
		t.Fatal(err)
	}

	v2, err := g.CreateVertex("2")
	if err != nil {
		t.Fatal(err)
	}
	sv2, err := EncodeVertexJSON(v2)
	if err != nil {
		t.Fatal(err)
	}

	e, err := g.CreateEdge("1", "testlabel", v1.Key(), v2.Key())
	if err != nil {
		t.Fatal(err)
	}
	se, err := EncodeEdgeJSON(e)
	if err != nil {
		t.Fatal(err)
	}

	g = NewSimpleGraph()
	_, err = DecodeVertexJSON(g, sv1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodeVertexJSON(g, sv2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodeEdgeJSON(g, se)
	if err != nil {
		t.Fatal(err)
	}

	e, err = g.GetEdge("1")
	if err != nil {
		t.Fatal(err)
	}

	if e.Head().Key() != v1.Key() {
		t.Fatal("Could not retrieve head")
	}

	if e.Tail().Key() != v2.Key() {
		t.Fatal("Could not retrieve tail")
	}

	if v1.Incoming()[0].Key() != e.Key() {
		t.Fatal("Could not retrieve edge")
		return
	}

	if v2.Outgoing()[0].Key() != e.Key() {
		t.Fatal("Could not retrieve edge")
	}

}
