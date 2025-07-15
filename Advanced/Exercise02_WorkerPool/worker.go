package exercise02workerpool

import "time"

// Worker represents a single processing unit in the worker pool.
// Its responsibility is to take tasks from an input channel, process them,
// and then send the results to an output channel.
type Worker struct {
	ID            int         // Unique identifier for the worker, useful for logging and debugging.
	TaskChannel   <-chan Task // A receive-only channel from which the worker receives tasks.
	ResultChannel chan<- Task // A send-only channel to which the worker sends processed tasks (results).
}

// NewWorker creates and returns a new instance of a Worker.
// It initializes the worker with an ID and the channels it will use for communication.
func NewWorker(id int, taskChan <-chan Task, resultChannel chan<- Task) *Worker {
	return &Worker{
		ID:            id,            // Assigns the given ID to the worker.
		TaskChannel:   taskChan,      // Assigns the task input channel.
		ResultChannel: resultChannel, // Assigns the result output channel.
	}
}

// Start begins the worker's main processing loop.
// This method is designed to be run in its own goroutine.
func (w *Worker) Start() {
	// The 'for range' loop over a channel will continuously receive values
	// until the channel is closed. Once the channel is closed and all
	// values have been received, the loop will terminate.
	for task := range w.TaskChannel {
		// Simulate Processing time based on the task's defined complexity.
		// This uses time.Sleep to block the goroutine for a specified duration,
		// mimicking actual work being done that consumes time.
		time.Sleep(task.Complexity)

		// Perform the CPU-intensive calculation for the task.
		// The 'isPrime' function is called with the task's data.
		// The result of this computation is assigned to the 'Result' field of the task.
		// Since 'task' is a value received from a channel, modifying it here is safe
		// as it's a local copy, not shared with other goroutines concurrently.
		task.Result = isPrime(task.Data)
		// Set any error to nil, assuming successful processing for this example.
		task.Err = nil

		// Send the processed task (now containing the result) back to the ResultChannel.
		// This sends the task to the consumer or further processing stages.
		w.ResultChannel <- task
	}
	// The loop exits when w.TaskChannel is closed by the Producer.
	// At this point, the worker goroutine will finish its execution.
}
