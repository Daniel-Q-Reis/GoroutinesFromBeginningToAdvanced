package exercise02workerpool

import "sync"

// Pool manages the worker pool and communication channels
type Pool struct {
	TaskChan    chan Task
	ResultChan  chan Task
	workerCount int
}

// NewPool creates a new worker pool instance
func NewPool(workerCount int) *Pool {
	return &Pool{
		TaskChan:    make(chan Task),                //Unbuffered for backpressure
		ResultChan:  make(chan Task, workerCount*2), //Buffered for throughput
		workerCount: workerCount,
	}
}

// Start launches all workers in the pool
func (p *Pool) Start() {
	var wg sync.WaitGroup

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		worker := NewWorker(i, p.TaskChan, p.ResultChan)
		go func() {
			defer wg.Done()
			worker.Start() //Start each worker in its own goroutine
		}()
	}

	go func() {
		wg.Wait()
		close(p.ResultChan)
	}()
}
