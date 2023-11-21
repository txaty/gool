# Gool

[![Go Reference](https://pkg.go.dev/badge/github.com/txaty/gool.svg)](https://pkg.go.dev/github.com/txaty/gool)
[![Go Report Card](https://goreportcard.com/badge/github.com/txaty/gool)](https://goreportcard.com/report/github.com/txaty/gool)
[![codecov](https://codecov.io/gh/txaty/gool/branch/main/graph/badge.svg?token=M02CIBSXFR)](https://codecov.io/gh/txaty/gool)

A generic goroutine pool just like Python ThreadPoolExecutor.

Gool provides the following methods:

- ```Submit```: Submit a task and return the result (if any).
- ```AsyncSubmit```: Submit a task and return a future of the result (if any), the future is the result
  channel.
- ```Map```: Submit a bundle of tasks and return the results in order (if any).
- ```AsyncMap```: Submit a bundle of tasks and return the futures of the results (if any), the futures are the result
  channels.

To use Gool, you need to define the following:

- Handler function: ```handler func(A) R```, and
- Argument: ```arg A```.

With types ```A``` and ```R``` being arbitrary types.

You can also specify the number of workers ```numWorkers``` and the task queue size ```cap``` when creating a new pool.
