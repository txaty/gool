package gool

import (
	"context"
)

type Pool struct {
	lmt        int
	workerPool chan *worker
	numWorkers int
	ctx        context.Context
}

func NewPool(lmt int) *Pool {
	p := &Pool{
		lmt:        lmt,
		workerPool: make(chan *worker, lmt),
		ctx:        context.Background(),
	}
	return p
}

func (p *Pool) Submit(handler func(interface{}), args interface{}) {
	w := p.getWorker()
	w.task <- task{handler, args}
	p.workerPool <- w
}

func (p *Pool) getWorker() *worker {
	for {
		select {
		case w := <-p.workerPool:
			return w
		default:
			if p.numWorkers < p.lmt {
				p.numWorkers++
				return newWorker(p.ctx)
			}
		}
	}
}

func (p *Pool) Close() {
	p.ctx.Done()
}
