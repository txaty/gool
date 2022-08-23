package gool

type Pool struct {
	numWorkers int
	jobChan    chan Task
}

// NewPool creates a new goroutine pool with the given number of workers and job queue capacity.
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

// Submit submits a task and waits for the result
func (p *Pool) Submit(handler func(interface{}) interface{}, args interface{}) interface{} {
	result := make(chan interface{})
	p.jobChan <- Task{
		handler: handler,
		args:    args,
		result:  result,
	}
	return <-result
}

// SubmitBatch submits a batch of tasks and waits for the results
func (p *Pool) SubmitBatch(handler func(interface{}) interface{},
	args []interface{}) []interface{} {
	resultChanList := make([]chan interface{}, len(args))
	for i := 0; i < len(args); i++ {
		resultChanList[i] = make(chan interface{})
		p.jobChan <- Task{
			handler: handler,
			args:    args[i],
			result:  resultChanList[i],
		}
	}
	results := make([]interface{}, len(args))
	for i := 0; i < len(args); i++ {
		results[i] = <-resultChanList[i]
	}
	return results
}

// Close closes the pool and waits for all the workers to stop
func (p *Pool) Close() {
	for i := 0; i < p.numWorkers; i++ {
		p.jobChan <- Task{
			stop: true,
		}
	}
}
