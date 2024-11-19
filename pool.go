package gpool

import (
	"sync"
)

type Worker interface {
	Do(i interface{})
}

type processor struct{}

type Pool struct {
	mu          *sync.Mutex
	max         int
	wg          sync.WaitGroup
	chProcessor chan *processor
}

func Init(max int) *Pool {
	if max <= 0 {
		max = 5
	}
	p := &Pool{
		mu:          &sync.Mutex{},
		max:         max,
		wg:          sync.WaitGroup{},
		chProcessor: make(chan *processor, max),
	}

	for i := 0; i < max; i++ {
		p.chProcessor <- &processor{}
	}

	return p
}

func (p *Pool) Do(w Worker, item interface{}) {
	//do sth
	p.wg.Add(1)
	go func() {
		//consume
		processor := <-p.chProcessor
		w.Do(item)
		p.wg.Done()
		//send back
		p.chProcessor <- processor
	}()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
