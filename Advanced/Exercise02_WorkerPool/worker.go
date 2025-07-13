package exercise02workerpool

import "time"

// Worker processes tasks from input channel and sends to output channel
type Worker struct {
	ID            int
	TaskChannel   <-chan Task
	ResultChannel chan<- Task
}

// NewWorker creates a new Worker instance
func NewWorker(id int, taskChan <-chan Task, resultChannel chan<- Task) *Worker {
	return &Worker{
		ID:            id,
		TaskChannel:   taskChan,
		ResultChannel: resultChannel,
	}
}

// Start begins the worker's processing loop
func (w *Worker) Start() {
	for task := range w.TaskChannel {
		// Simulate Processing time
		time.Sleep(task.Complexity)
		//Process the task
		task.Result = isPrime(task.Data) // Performa the CPU-Intensive calculation
		task.Err = nil

		w.ResultChannel <- task
	}
}
