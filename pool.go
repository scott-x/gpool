package gpool

import (
	"log"
	"sync"
	"time"
)

type Worker interface {
	Do(item interface{})
}

type processor struct {
	index int
}

type Pool struct {
	mu              *sync.RWMutex
	name            string
	max             int
	workingNum      int
	chProcessor     chan *processor
	isLoggerEnabled bool
}

func New(name string, max int) *Pool {
	if max <= 0 {
		max = 5
	}
	p := &Pool{
		mu:              &sync.RWMutex{},
		name:            name,
		max:             max,
		chProcessor:     make(chan *processor, max),
		isLoggerEnabled: false,
		workingNum:      0,
	}

	for i := 0; i < max; i++ {
		p.chProcessor <- &processor{
			index: i,
		}
	}

	return p
}

func (p *Pool) EnableLogger() *Pool {
	p.isLoggerEnabled = true
	return p
}

func (p *Pool) Treat(w Worker, item interface{}) {
	//do sth
	//consume block
	processor := <-p.chProcessor

	p.mu.Lock()
	p.workingNum++
	p.mu.Unlock()

	if p.isLoggerEnabled {
		log.Printf("%s_proceessor_%d: %v\n", p.name, processor.index, item)
	}
	w.Do(item)

	//send back
	p.chProcessor <- processor

	p.mu.Lock()
	p.workingNum--
	p.mu.Unlock()
}

func (p *Pool) Wait() {
	for {
		var num int

		time.Sleep(time.Second)
		p.mu.Lock()
		num = p.workingNum
		p.mu.Unlock()

		if num == 0 {
			return
		}
	}
}
