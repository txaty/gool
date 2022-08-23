package gool

type Task struct {
	handler func(interface{})
	args    interface{}
	finish  chan struct{}
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
	for {
		select {
		case job := <-w.jobChan:
			if job.stop {
				return
			}
			job.handler(job.args)
			job.finish <- struct{}{}
		}
	}
}
