package transactionedge

import (
	"github.com/joernweissenborn/propertygraph2go/propertygraph"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionticket"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionitem"
)


func New(original propertygraph.Edge,ticket *transactionticket.TransactionTicket) *TransactionEdge{
	tv := new(TransactionEdge)
	tv.original = original
	tv.TransactionItem = transactionitem.New(tv.original.Id(), ticket, original.Properties().(*transactionitem.Properties) )
	return tv
}

func Create(id string,label string, head propertygraph.Vertex, tail propertygraph.Vertex, properties interface{},ticket *transactionticket.TransactionTicket) *TransactionEdge{
	te := new(TransactionEdge)
	te.label = label
	te.head = head
	te.tail = tail
	te.TransactionItem = transactionitem.Create(id, ticket, properties)
	return te
}

type TransactionEdge struct {
	*transactionitem.TransactionItem
	original propertygraph.Edge
	label string
	head propertygraph.Vertex
	tail propertygraph.Vertex
}

func (te TransactionEdge) Label()      string{
	te.DirtyCheck()

	if te.Written {
		return te.label
	}
	return te.original.Label()
}
func (te TransactionEdge) Head() propertygraph.Vertex {
	te.DirtyCheck()
	if te.Written {
		return te.head
	}
	return te.original.Head()
}
func (te TransactionEdge) Tail() propertygraph.Vertex {
	te.DirtyCheck()
	if te.Written {
		return te.tail
	}
	return te.original.Tail()
}

