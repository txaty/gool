package gool

// Pool implements a simple goroutine pool
type Pool[A, R any] struct {
	numWorkers int
	jobChan    chan Task[A, R]
}

// NewPool creates a new goroutine pool with the given number of workers and job queue capacity.
func NewPool[A, R any](numWorkers, cap int) *Pool[A, R] {
	p := &Pool[A, R]{
		numWorkers: numWorkers,
		jobChan:    make(chan Task[A, R], cap),
	}
	for i := 0; i < numWorkers; i++ {
		newWorker(p.jobChan)
	}
	return p
}

// Submit submits a task and waits for the result
func (p *Pool[A, R]) Submit(handler func(A) R, args A) R {
	result := make(chan R)
	p.jobChan <- Task[A, R]{
		handler: handler,
		args:    args,
		result:  result,
	}
	return <-result
}

// Map submits a batch of tasks and waits for the results
func (p *Pool[A, R]) Map(handler func(A) R, args []A) []R {
	resultChanList := make([]chan R, len(args))
	for i := 0; i < len(args); i++ {
		resultChanList[i] = make(chan R)
		p.jobChan <- Task[A, R]{
			handler: handler,
			args:    args[i],
			result:  resultChanList[i],
		}
	}
	results := make([]R, len(args))
	for i := 0; i < len(args); i++ {
		results[i] = <-resultChanList[i]
	}
	return results
}

// Close closes the pool and waits for all the workers to stop
func (p *Pool[A, R]) Close() {
	for i := 0; i < p.numWorkers; i++ {
		p.jobChan <- Task[A, R]{
			stop: true,
		}
	}
}
