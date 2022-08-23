package gool

// Task is a job to be executed by a worker
type Task struct {
	handler func(interface{}) interface{}
	args    interface{}
	result  chan interface{}
	stop    bool
}

type worker struct {
	jobChan chan Task
}

func newWorker(jobChan chan Task) *worker {
	w := &worker{
		jobChan: jobChan,
	}
	go w.run()
	return w
}

func (w *worker) run() {
	for job := range w.jobChan {
		if job.stop {
			return
		}
		job.result <- job.handler(job.args)
	}
}
