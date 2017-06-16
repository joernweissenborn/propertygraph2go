package propertygraph2go

import "fmt"

type KeyGen struct {
	next_vertex int
	next_edge   int
}

func NewKeyGen() (k *KeyGen) { return &KeyGen{} }

func (k *KeyGen) NextVertex() (i int) {
	i = k.next_vertex
	k.next_vertex++
	return
}

func (k *KeyGen) NextEdge() (i int) {
	i = k.next_edge
	k.next_edge++
	return
}

type StringKeyGen struct {
	next_vertex int
	next_edge   int
	prefix      string
}

func NewStringGen(prefix string) (k *StringKeyGen) { return &StringKeyGen{prefix: prefix} }

func (k *StringKeyGen) NextVertex() (i string) {
	i = fmt.Sprintf("%s%d", k.prefix, k.next_vertex)
	k.next_vertex++
	return
}

func (k *StringKeyGen) NextEdge() (i string) {
	i = fmt.Sprintf("%s%d", k.prefix, k.next_edge)
	k.next_edge++
	return
}
