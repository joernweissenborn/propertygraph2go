package PropertyGraph

import "testing"

func TestReadEmpty(t *testing.T) {
	graph := New()
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
		t.Error("Empty edge query did not return nil")
	}

}

func TestWriteReadVertex(t *testing.T) {
	graph := New()
	type testprops struct {
		Name string
	}
	var props testprops
	props.Name = "Habib"

	graph.CreateVertex("test", props)

	v := graph.GetVertex("test")

	rprops, found := v.Properties.(testprops)

	if !found {
		t.Error("Write Read test type assertion failed")
		return
	}

	if rprops.Name != "Habib" {
		t.Error("Write Read test failed. Expected 'Habib' but got", rprops.Name)
		return
	}

	t.Log("Write Read test passed...")

}

func TestWriteReadEdge(t *testing.T) {
	graph := New()

	v1 := graph.CreateVertex("1", nil)
	v2 := graph.CreateVertex("2", nil)

	graph.CreateEdge("1", "testlabel", v1, v2, nil)

	e := graph.GetEdge("1")

	if e == nil {
		t.Error("Could not retrieve edge")
		return
	}

	if e.Head.Id != v1.Id {
		t.Error("Could not retrieve head")
		return
	}

	if e.Tail.Id != v2.Id {
		t.Error("Could not retrieve tail")
		return
	}

	if v1.Incoming[0].Tail.Id != v2.Id {
		t.Error("Could not retrieve tail")
		return
	}

	if v2.Outgoing[0].Head.Id != v1.Id {
		t.Error("Could not retrieve head")
		return
	}

}
