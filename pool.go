package gpool

import (
	"sync"
)

type Worker interface {
	Do(i interface{})
}

type Processor struct {
}

type Pool struct {
	mu *sync.Mutex
	// Processors  []*Processor
	max         int
	wg          sync.WaitGroup
	ChProcessor chan *Processor
}

func Init(max int) *Pool {
	if max <= 0 {
		max = 5
	}
	p := &Pool{
		mu:          &sync.Mutex{},
		max:         max,
		wg:          sync.WaitGroup{},
		ChProcessor: make(chan *Processor, max),
	}

	for i := 0; i < max; i++ {
		p.ChProcessor <- &Processor{}
	}

	return p
}

func (p *Pool) Do(w Worker, item interface{}) {
	//do sth
	p.wg.Add(1)
	go func() {
		//consume
		processor := <-p.ChProcessor
		w.Do(item)
		p.wg.Done()
		//send back
		p.ChProcessor <- processor
	}()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
