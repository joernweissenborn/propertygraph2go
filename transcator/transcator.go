package transcator

import (
	"github.com/joernweissenborn/propertygraph2go/propertygraph"
	"github.com/joernweissenborn/stream2go"
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionticket"
)

type Transcator struct {
	g propertygraph.PropertyGraph
	nextTicketNr int64
	transactionstream stream2go.StreamController
	commitchan chan *transactionticket.TransactionTicket
}

func New(g propertygraph.PropertyGraph) (ta Transcator){
	ta.transactionstream = stream2go.New()
	ta.commitchan = make(chan *transactionticket.TransactionTicket)
	go ta.commitRunner()
	ta.g = g
	return
}


func (t *Transcator) Transaction(tf TransactionFunc){

	ticket := transactionticket.New(t.nextTicketNr)
	t.nextTicketNr++
	ticket.Complete.Then(t.commit)
	ta := NewTransaction(t.g,ticket,tf)
	transact(ta)

}

func transact(ta *Transaction)  {

	f := false
	n := 0
	for !f {
		ta.transactions[n](ta)
		n++
		f = len(ta.transactions) == n
	}
	ta.ticket.CompleteTransaction()
	ta.ticket.Release.WaitUntilComplete()
	return
}



func (t *Transcator) commit(ticket interface {}) interface {} {
	t.commitchan <- ticket.(*transactionticket.TransactionTicket)
	return nil
}


func (t Transcator) commitRunner()  {
	finish := make(chan struct {})
	tickets := []*transactionticket.TransactionTicket{}
	var curTicketNumber int64 = 0
	for {
		select {
		case <-finish:
			curTicketNumber++
			for _,ticket := range tickets {
				if ticket.Number == curTicketNumber {
					ticket.Commit.Complete(finish)
					break
				}
			}
		case ticket := <-t.commitchan:
			if ticket.Number == curTicketNumber {
				ticket.Commit.Complete(finish)
			} else{
				tickets = append(tickets,ticket)
			}
		}
	}

	return
}

