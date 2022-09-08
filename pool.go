// MIT License
//
// Copyright (c) 2022 Tommy TIAN
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gool

import (
	"runtime"
)

// Pool implements a simple goroutine pool.
type Pool[A, R any] struct {
	numWorkers int
	taskChan   chan task[A, R]
}

// NewPool creates a new goroutine pool with the given number of workers and job queue capacity.
// If numWorkers is less than 1, it will be set to the number of CPUs.
// If cap (task queue capacity) is less than 1, it will be set to twice the number of workers.
func NewPool[A, R any](numWorkers, cap int) *Pool[A, R] {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	if cap <= 0 {
		cap = 2 * numWorkers
	}
	p := &Pool[A, R]{
		numWorkers: numWorkers,
		taskChan:   make(chan task[A, R], cap),
	}
	for i := 0; i < numWorkers; i++ {
		newWorker(p.taskChan)
	}
	return p
}

// Submit submits a task and waits for the result.
func (p *Pool[A, R]) Submit(handler func(A) R, args A) R {
	result := p.AsyncSubmit(handler, args)
	return <-result
}

// AsyncSubmit submits a task and returns the channel to wait for the result.
func (p *Pool[A, R]) AsyncSubmit(handler func(A) R, args A) chan R {
	resChan := make(chan R)
	p.taskChan <- task[A, R]{
		handler: handler,
		args:    args,
		result:  resChan,
	}
	return resChan
}

// Map submits a batch of tasks and waits for the results.
func (p *Pool[A, R]) Map(handler func(A) R, args []A) []R {
	resultChanList := p.AsyncMap(handler, args)
	results := make([]R, len(args))
	for i := 0; i < len(args); i++ {
		results[i] = <-resultChanList[i]
	}
	return results
}

// AsyncMap submits a batch of tasks and returns the channel to wait for the results.
func (p *Pool[A, R]) AsyncMap(handler func(A) R, args []A) []chan R {
	resultChanList := make([]chan R, len(args))
	for i := 0; i < len(args); i++ {
		resultChanList[i] = make(chan R)
		p.taskChan <- task[A, R]{
			handler: handler,
			args:    args[i],
			result:  resultChanList[i],
		}
	}
	return resultChanList
}

// Close closes the pool and waits for all the workers to stop.
func (p *Pool[A, R]) Close() {
	for i := 0; i < p.numWorkers; i++ {
		p.taskChan <- task[A, R]{
			stop: true,
		}
	}
}
