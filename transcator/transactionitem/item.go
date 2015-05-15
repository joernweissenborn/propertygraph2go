package transactionitem

import (
	"github.com/joernweissenborn/propertygraph2go/transcator/transactionticket"
	"sync"
	"github.com/joernweissenborn/future2go"
)

func New(id string, ticket *transactionticket.TransactionTicket, props *Properties) *TransactionItem{
	ti := new(TransactionItem)
	ti.id = id
	ti.ticket = ticket
	ti.props = props
	ti.lockRead()
	return ti
}

func Create(id string, ticket *transactionticket.TransactionTicket, properties interface {}) *TransactionItem{
	ti := new(TransactionItem)
	ti.id = id
	ti.ticket = ticket
	ti.props = &Properties{new(sync.Mutex),properties,[]*transactionticket.TransactionTicket{}, nil}
	ti.Written = true

	return ti
}

type TransactionItem struct {
	ticket *transactionticket.TransactionTicket
	id string
	props *Properties
	Written bool
}

func (ti TransactionItem) DirtyCheck() {
	ti.ticket.DirtyCheck()
}
func (ti TransactionItem) Id() string {
	return ti.id
}

func (ti TransactionItem) Properties() interface{} {
	return ti.props.originalProperties
}

func (ti TransactionItem) Props() *Properties {
	return ti.props
}

func (tv TransactionItem) Lock(){
	p := tv.props
	p.Lock()
	defer p.Unlock()


	if !p.writeLockClear() {
		if p.haveWriteLock(tv.ticket){
			return
		}

		if p.haveHigherLock(tv.ticket) {
			p.writeLock.Dirty.Complete(tv.ticket)
		} else {
			c := make(chan struct{})

			p.writeLock.Release.Then(onrelease(c))
			p.writeLock.Dirty.Then(onnewlock(c))

			<-c

		}
	}

	var wg sync.WaitGroup
	for _ , t := range p.readLocks {
		if t.Number < tv.ticket.Number {
			wg.Add(1)
			release := func(interface{}) interface{} {
				wg.Done()
				return nil
			}
			t.Release.Then(release)
			t.Dirty.Then(func(d interface{}) interface{} {
				d.(transactionticket.TransactionTicket).Release.Then(release)
				return nil
			})
		} else {
			p.writeLock.Dirty.Complete(tv.ticket)
		}
	}

	wg.Wait()

	p.writeLock = tv.ticket

}


func (tv TransactionItem)  lockRead(){
	p := tv.props
	p.Lock()
	defer p.Unlock()

	if p.haveReadLock(tv.ticket) || p.haveWriteLock(tv.ticket){

		return
	}

	if !p.writeLockClear() {
		c := make(chan struct{})
		p.writeLock.Release.Then(onrelease(c))
		p.writeLock.Dirty.Then(onnewlock(c))
		p.Unlock()
		<-c
		p.Lock()
	}

	p.readLocks = append(p.readLocks,tv.ticket)
}

func onrelease(c chan struct{}) future2go.CompletionFunc {
	return func(interface{}) interface{} {
		c <- struct{}{}
		return nil
	}
}
func onnewlock(c chan struct{}) future2go.CompletionFunc {
	return func(newlock interface{}) interface{} {
		newlock.(*transactionticket.TransactionTicket).Release.Then(onrelease(c))
		newlock.(*transactionticket.TransactionTicket).Dirty.Then(onnewlock(c))
		return nil
	}
}

func (p Properties) writeLockClear() bool{
	return p.writeLock == nil
}

func (p Properties) haveWriteLock(ticket *transactionticket.TransactionTicket) bool {
	if (p.writeLock == nil) {return false}
	return p.writeLock.Number == ticket.Number
}
func (p Properties) haveHigherLock(ticket *transactionticket.TransactionTicket) bool {
	return p.writeLock.Number > ticket.Number
}

func (p Properties) haveReadLock(ticket *transactionticket.TransactionTicket) bool{

	for _,t := range p.readLocks {
		if t.Number == ticket.Number{
			return true
		}
	}

	return false
}

type Properties struct {
	*sync.Mutex
	originalProperties interface {}
	readLocks []*transactionticket.TransactionTicket
	writeLock *transactionticket.TransactionTicket
}
