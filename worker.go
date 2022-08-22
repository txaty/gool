package gool

import (
	"context"
)

type task struct {
	handler func(interface{})
	args    interface{}
}

type worker struct {
	task chan task
	ctx  context.Context
}

func newWorker(ctx context.Context) *worker {
	w := &worker{
		task: make(chan task),
		ctx:  ctx,
	}
	go w.run()
	return w
}

func (w *worker) run() {
	for {
		select {
		case t := <-w.task:
			t.handler(t.args)
		case <-w.ctx.Done():
			return
		}
	}
}
