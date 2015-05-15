package transactionticket

import "github.com/joernweissenborn/future2go"

func New(Number int64) *TransactionTicket {
	return &TransactionTicket{Number, future2go.New(), future2go.New(), future2go.New(), future2go.New()}
}

type TransactionTicket struct {
	Number int64
	Complete *future2go.Future
	Commit *future2go.Future
	Release *future2go.Future
	Dirty *future2go.Future
}


func (t *TransactionTicket) CompleteTransaction(){
	t.Complete.Complete(t)
}

func (t *TransactionTicket) DirtyCheck(){
	if t.Dirty.IsComplete() {
		panic("Dirty")
	}
}

