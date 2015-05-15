package transcator

import (
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionticket"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionvertex"
	"github.com/joernweissenborn/propertygraph2go/propertygraph"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionedge"

)

func NewTransaction(original propertygraph.PropertyGraph, ticket *transactionticket.TransactionTicket, tf TransactionFunc) *Transaction{
	ta := new(Transaction)
	ta.original = original
	ticket.Commit.Then(ta.commit)
	ta.ticket = ticket
	ta.transactions = []TransactionFunc{tf}
	return ta
}

type Transaction struct {

	transactions []TransactionFunc

	original propertygraph.PropertyGraph

	ticket *transactionticket.TransactionTicket

	createdVertices []propertygraph.Vertex
	touchedVertices []propertygraph.Vertex
	removedVertices []string

	createdEdges []propertygraph.Edge
	touchedEdges []propertygraph.Edge
	removedEdges []string
}

func (t *Transaction) commit(c interface {})interface {}{

	for _, v := range t.removedVertices {
		t.original.RemoveVertex(v)
	}
	for _, e := range t.removedEdges {
		if t.original.GetEdge(e) != nil {
			t.original.RemoveEdge(e)
		}
	}
	for _, v := range t.createdVertices {
		t.original.CreateVertex(v.Id(),v.(*transactionvertex.TransactionVertex).Props())
	}
	for _, e := range t.createdEdges{
		hv := t.original.GetVertex(e.Head().Id())
		tv := t.original.GetVertex(e.Tail().Id())
		t.original.CreateEdge(e.Id(),e.Label(),hv,tv,e.(*transactionedge.TransactionEdge).Props())
	}
	t.ticket.Release.Complete(nil)
	c.(chan struct{}) <- struct{}{}
	return nil
}

func (t *Transaction) DirtyCheck() {
	t.ticket.DirtyCheck()
}
func (t *Transaction) CreateVertex(id string, properties interface{}) propertygraph.Vertex {
	t.DirtyCheck()
	tv := transactionvertex.Create(id,properties, t.ticket)

	t.createdVertices = append(t.createdVertices,tv)

	return tv
}

func (t *Transaction) RemoveVertex(id string) {
	t.DirtyCheck()
	t.removedVertices = append(t.removedVertices,id)
}

func (t Transaction) isRemovedVertex(id string) bool{
	for _, rv := range t.removedVertices {
		if rv == id {
			return true
		}
	}
	return false
}
func (t Transaction) isCreatedVertex(id string) (tv propertygraph.Vertex, is bool){
	for _, cv := range t.createdVertices{
		if cv.Id() == id {
			return cv, true
		}
	}
	return nil, false
}
func (t Transaction) isTouchedVertex(id string) (tv propertygraph.Vertex, is bool){
	for _, v := range t.touchedVertices{
		if v.Id() == id {
			return v, true
		}
	}
	return nil, false
}

func (t Transaction) GetVertex(id string) (v propertygraph.Vertex) {
	t.DirtyCheck()
	if t.isRemovedVertex(id) {
		return nil
	}
	v,f := t.isCreatedVertex(id)
	if f {return}

	v,f = t.isTouchedVertex(id)
	if f {return}

	o := t.original.GetVertex(id)
	if o != nil {
		v = transactionvertex.New(o,t.ticket)
		t.touchedVertices = append(t.touchedVertices)
	}
	return
}

func (t *Transaction) CreateEdge(id string, label string, head propertygraph.Vertex, tail propertygraph.Vertex, properties interface{}) propertygraph.Edge {
	t.DirtyCheck()

	te := transactionedge.Create(id,label,head,tail,properties, t.ticket)

	head.AddIncomingEdge(te)
	tail.AddOutgoingEdge(te)

	t.createdEdges = append(t.createdEdges,te)

	return te
}
func (t Transaction) RemoveEdge(id string) {
	t.DirtyCheck()
	t.removedEdges = append(t.removedEdges,id)
}

func (t Transaction) isRemovedEdge(id string) bool{
	for _, rv := range t.removedEdges{
		if rv == id {
			return true
		}
	}
	return false
}
func (t Transaction) isCreatedEdge(id string) (te propertygraph.Edge, is bool){
	for _, ce := range t.createdEdges{
		if ce.Id() == id {
			return ce, true
		}
	}
	return nil, false
}
func (t Transaction) isTouchedEdge(id string) (te propertygraph.Edge, is bool){
	for _, e := range t.touchedEdges{
		if e.Id() == id {
			return e, true
		}
	}
	return nil, false
}
func (t Transaction) GetEdge(id string) (e propertygraph.Edge) {
	t.DirtyCheck()
	if t.isRemovedEdge(id) {
		return nil
	}
	e,f := t.isCreatedEdge(id)
	if f {return}

	e,f = t.isTouchedEdge(id)
	if f {return}

	o := t.original.GetEdge(id)
	if o != nil {
		e = transactionedge.New(o,t.ticket)
		t.touchedEdges = append(t.touchedEdges,e)
	}
	return
}
