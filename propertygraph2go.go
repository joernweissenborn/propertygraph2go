package propertygraph2go

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Key interface{}

var ErrKeyNotFound = errors.New("Key not found")

type VertexWalker func(v Vertex)
type EdgeWalker func(e Edge)

type Graph interface {
	CreateVertex(key Key) (v Vertex, err error)
	GetVertex(key Key) (v Vertex, err error)
	WalkVertices(VertexWalker) (err error)
	RemoveVertex(key Key) (err error)

	CreateEdge(k Key, label string, head, tail Key) (e Edge, err error)
	GetEdge(k Key) (e Edge, err error)
	WalkEdges(EdgeWalker) (err error)
	RemoveEdge(k Key) (err error)

	WithProperties
}

type Edge interface {
	Key() (k Key)
	Label() (label string)
	Head() (v Vertex)
	Tail() (v Vertex)
	WithProperties
}

type Vertex interface {
	Key() (k Key)
	Incoming() (incoming []Edge)
	Outgoing() (Outgoing []Edge)
	WithProperties
}

func StringVertex(v Vertex) (s string) {
	s = fmt.Sprintf("Vertex (%s)", v.Key())
	s = fmt.Sprintf("%s\n\tProperties:", s)
	props, _ := v.Properties()
	for n, v := range props {
		s = fmt.Sprintf("%s {%s: %v}", s, n, v)
	}

	s = fmt.Sprintf("%s\n\tIncoming:", s)
	for _, e := range v.Incoming() {
		s = fmt.Sprintf("%s {%s %s: %s -> %s}", s, e.Label(), e.Key(), e.Tail().Key(), e.Head().Key())
	}

	s = fmt.Sprintf("%s\n\tOutgoing:", s)
	for _, e := range v.Outgoing() {
		s = fmt.Sprintf("%s {%s %s: %s -> %s}", s, e.Label(), e.Key(), e.Tail().Key(), e.Head().Key())
	}
	return
}

func PrintVertex(v Vertex) {
	fmt.Println(StringVertex(v))
}

func PrintVertices(vs ...Vertex) {
	for _, v := range vs {
		PrintVertex(v)
	}
}

func PrintAllVertices(g Graph) {
	g.WalkVertices(PrintVertex)
}

type EncodingFunc func(interface{}) ([]byte, error)
type DecodingFunc func([]byte, interface{}) error

type EncodedVertex struct {
	Key        Key
	Properties map[string]interface{}
}

type EncodedEdge struct {
	Key        Key
	Label      string
	Properties map[string]interface{}
	Head       Key
	Tail       Key
}

func EncodeVertexJSON(e Vertex) (data []byte, err error) { return EncodeVertex(json.Marshal, e) }
func EncodeVertex(enc EncodingFunc, v Vertex) (data []byte, err error) {
	p, err := v.Properties()
	if err != nil {
		return
	}
	ev := EncodedVertex{
		Key:        v.Key(),
		Properties: p,
	}
	return enc(ev)
}

func DecodeVertexJSON(g Graph, data []byte) (v Vertex, err error) {
	return DecodeVertex(json.Unmarshal, g, data)
}
func DecodeVertex(dec DecodingFunc, g Graph, data []byte) (v Vertex, err error) {
	var ev EncodedVertex
	err = dec(data, &ev)
	if err != nil {
		return
	}
	v, err = g.CreateVertex(ev.Key)
	if err != nil {
		return
	}
	for k, val := range ev.Properties {
		err = v.SetProperty(k, val)
		if err != nil {
			return
		}
	}
	return
}

func EncodeEdgeJSON(e Edge) (data []byte, err error) { return EncodeEdge(json.Marshal, e) }
func EncodeEdge(enc EncodingFunc, e Edge) (data []byte, err error) {
	p, err := e.Properties()
	if err != nil {
		return
	}
	ee := EncodedEdge{
		Key:        e.Key(),
		Label:      e.Label(),
		Properties: p,
		Head:       e.Head().Key(),
		Tail:       e.Tail().Key(),
	}
	return enc(ee)
}

func DecodeEdgeJSON(g Graph, data []byte) (e Edge, err error) {
	return DecodeEdge(json.Unmarshal, g, data)
}
func DecodeEdge(dec DecodingFunc, g Graph, data []byte) (e Edge, err error) {
	var ee EncodedEdge
	err = dec(data, &ee)
	if err != nil {
		return
	}
	e, err = g.CreateEdge(ee.Key, ee.Label, ee.Head, ee.Tail)
	if err != nil {
		return
	}
	for k, v := range ee.Properties {
		err = e.SetProperty(k, v)
		if err != nil {
			return
		}
	}
	return
}

type KeyToPath func(Key) string

func EncodeGraph(enc EncodingFunc, g Graph) (vs map[Key][]byte, es map[Key][]byte) {
	vs = map[Key][]byte{}
	walkVertices := func(v Vertex) {
		sv, err := EncodeVertex(enc, v)
		if err != nil {
			panic(err)
		}
		vs[v.Key()] = sv
	}
	es = map[Key][]byte{}
	walkEdges := func(e Edge) {
		se, err := EncodeEdge(enc, e)
		if err != nil {
			panic(err)
		}
		es[e.Key()] = se
	}
	g.WalkVertices(walkVertices)
	g.WalkEdges(walkEdges)
	return
}
