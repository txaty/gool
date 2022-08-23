package gool

type Pool struct {
	numWorkers int
	jobChan    chan Task
}

func NewPool(numWorkers, cap int) *Pool {
	p := &Pool{
		numWorkers: numWorkers,
		jobChan:    make(chan Task, cap),
	}
	for i := 0; i < numWorkers; i++ {
		newWorker(p.jobChan)
	}
	return p
}

func (p *Pool) Submit(handler func(interface{}), args interface{}) {
	finish := make(chan struct{})
	p.jobChan <- Task{
		handler: handler,
		args:    args,
		finish:  finish,
	}
	<-finish
}

func (p *Pool) SubmitBatch(handler func(interface{}), args []interface{}) {
	finishes := make([]chan struct{}, len(args))
	for i := 0; i < len(args); i++ {
		finishes[i] = make(chan struct{})
		p.jobChan <- Task{
			handler: handler,
			args:    args[i],
			finish:  finishes[i],
		}
	}
	for i := 0; i < len(args); i++ {
		<-finishes[i]
	}
}

func (p *Pool) Close() {
	for i := 0; i < p.numWorkers; i++ {
		p.jobChan <- Task{
			stop: true,
		}
	}
}
