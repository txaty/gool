package gool

// Task is a job to be executed by a worker
type Task[A, R any] struct {
	handler func(A) R
	args    A
	result  chan R
	stop    bool
}

type worker[A, R any] struct {
	jobChan chan Task[A, R]
}

func newWorker[A, R any](jobChan chan Task[A, R]) *worker[A, R] {
	w := &worker[A, R]{
		jobChan: jobChan,
	}
	go w.run()
	return w
}

func (w *worker[A, R]) run() {
	for job := range w.jobChan {
		if job.stop {
			return
		}
		job.result <- job.handler(job.args)
	}
}
