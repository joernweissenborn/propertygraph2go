package transactionvertex

import (
	"github.com/joernweissenborn/propertygraph2go/propertygraph"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionticket"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionitem"
)


func New(original propertygraph.Vertex,ticket *transactionticket.TransactionTicket) *TransactionVertex{
	tv := new(TransactionVertex)
	tv.original = original
	tv.TransactionItem = transactionitem.New(tv.original.Id(), ticket, original.Properties().(*transactionitem.Properties))
	return tv
}

func Create(id string, properties interface{},ticket *transactionticket.TransactionTicket) *TransactionVertex{
	tv := new(TransactionVertex)
	tv.TransactionItem = transactionitem.Create(id, ticket, properties)
	return tv
}

type TransactionVertex struct {
	*transactionitem.TransactionItem
	original propertygraph.Vertex
	incoming []propertygraph.Edge
	outgoing []propertygraph.Edge
	}

func (tv TransactionVertex) Incoming() []propertygraph.Edge {
	tv.DirtyCheck()
	if tv.Written {
		return tv.incoming
	}
	return tv.original.Incoming()
}

func (tv TransactionVertex) Outgoing()[]propertygraph.Edge {
	tv.DirtyCheck()
	if tv.Written {
		return tv.outgoing
	}
	return tv.original.Outgoing()
}
func (tv TransactionVertex) RemoveIncomingEdge(id string) {
	tv.DirtyCheck()
	tv.copy()
	for i, e := range tv.incoming {
		if e.Id() == id {
			tv.incoming[i] = tv.incoming[len(tv.incoming)-1]
			tv.incoming = tv.incoming[:len(tv.incoming)-2]
			return
		}
	}
}

func (tv TransactionVertex) RemoveOutgoingEdge(id string){
	tv.DirtyCheck()
	tv.copy()

	for i, e := range tv.outgoing {
		if e.Id() == id {
			tv.outgoing[i] = tv.outgoing[len(tv.outgoing)-1]
			tv.outgoing = tv.outgoing[:len(tv.outgoing)-2]
			return
		}
	}}

func (tv TransactionVertex) AddIncomingEdge(e propertygraph.Edge){
	tv.DirtyCheck()
	tv.copy()
	tv.incoming = append(tv.incoming,e)
}

func (tv TransactionVertex) AddOutgoingEdge(e propertygraph.Edge){
	tv.DirtyCheck()
	tv.copy()
	tv.outgoing = append(tv.outgoing,e)
}

func (tv *TransactionVertex) copy(){

	if !tv.Written {
		tv.Written = true
		tv.incoming = tv.original.Incoming()
		tv.outgoing = tv.original.Outgoing()
		tv.Lock()
	}
}



