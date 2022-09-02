# Gool

A generic goroutine pool just like Python ThreadPoolExecutor.

Gool provides the following methods:

- ```Submit```: Submit a task and return the result (if any).
- ```AsyncSubmit```: Submit a task and return a future of the result (if any), the future is actually the result
  channel.
- ```Map```: Submit a bundle of tasks and return the results in order (if any).
- ```AsyncMap```: Submit a bundle of tasks and return the futures of the results (if any), the futures are the result
  channels.
