package exercise02workerpool

import "sync" // Package for synchronization primitives like WaitGroup.

// Pool manages the creation and orchestration of the worker goroutines
// and the communication channels between producers, workers, and consumers.
type Pool struct {
	TaskChan    chan Task // Channel for tasks to be sent to workers. Unbuffered for backpressure.
	ResultChan  chan Task // Channel for results to be sent from workers to consumers. Buffered for throughput.
	workerCount int       // The number of worker goroutines in this pool.
}

// NewPool creates and returns a new Pool instance.
// It initializes the task and result channels with appropriate buffering.
func NewPool(workerCount int) *Pool {
	return &Pool{
		// TaskChan is unbuffered (make(chan Task)). This means a sender (Producer)
		// will block until a receiver (Worker) is ready to take the task.
		// This provides a critical backpressure mechanism, preventing the producer
		// from generating tasks faster than workers can consume them, thus
		// avoiding unbounded memory usage for tasks awaiting processing.
		TaskChan: make(chan Task),

		// ResultChan is buffered (make(chan Task, workerCount*2)).
		// A buffered channel allows a sender (Worker) to send results without blocking
		// immediately, as long as the buffer is not full. This increases throughput
		// by decoupling workers from the consumer, allowing workers to continue
		// processing new tasks while the consumer might be temporarily busy.
		// The buffer size (workerCount*2) is a common heuristic, providing some
		// slack without consuming excessive memory.
		ResultChan: make(chan Task, workerCount*2),

		workerCount: workerCount, // Stores the number of workers this pool will manage.
	}
}

// Start launches all worker goroutines and manages the graceful closing of the ResultChan.
// This method sets up the core concurrency of the worker pool.
func (p *Pool) Start() {
	// A WaitGroup is used to wait for all worker goroutines to complete their work.
	// It's local to the Pool.Start method because it's only concerned with the workers
	// managed by this specific pool instance.
	var wg sync.WaitGroup

	// Loop to launch the specified number of worker goroutines.
	for i := 0; i < p.workerCount; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each worker about to be launched.

		// Create a new Worker instance for each goroutine.
		// Each worker receives its unique ID, the shared TaskChan, and the shared ResultChan.
		worker := NewWorker(i, p.TaskChan, p.ResultChan)

		// Launch the worker's processing loop in a new goroutine.
		go func() {
			// Defer wg.Done() ensures that the WaitGroup counter is decremented
			// when this goroutine finishes, regardless of how it exits (e.g., normally, panics).
			defer wg.Done()
			// Call the worker's Start method. This method will block and process tasks
			// until the TaskChan is closed by the producer and drained.
			worker.Start() // Start each worker in its own goroutine
		}()
	}

	// Launch a separate goroutine to manage the closing of the ResultChan.
	// This is crucial for a graceful shutdown, as the consumer's 'for range' loop
	// on ResultChan will only terminate when ResultChan is closed.
	go func() {
		// wg.Wait() blocks this goroutine until the WaitGroup counter becomes zero.
		// This means it waits until ALL worker goroutines launched by this pool
		// have completed their execution (signaled by defer wg.Done()).
		wg.Wait()
		// Once all workers are done, close the ResultChan.
		// This signals to the consumer (and any other goroutines reading from ResultChan)
		// that no more results will be sent, allowing their 'for range' loops to exit gracefully.
		close(p.ResultChan)
	}()
}
