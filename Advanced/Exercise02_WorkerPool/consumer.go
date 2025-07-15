package exercise02workerpool

import (
	"fmt" // Package for formatted I/O, used for printing results to the console.
)

// Consumer is responsible for receiving and processing the results from the worker pool.
// It reads processed tasks from the result channel and displays their outcome.
type Consumer struct {
	ResultChan <-chan Task // A receive-only channel from which the consumer receives processed tasks.
}

// NewConsumer creates and returns a new Consumer instance.
// It initializes the consumer with the channel from which it will receive results.
func NewConsumer(resultChan <-chan Task) *Consumer {
	return &Consumer{
		ResultChan: resultChan, // Assigns the result input channel.
	}
}

// Start begins the consumer's main loop for processing results.
// This method is designed to be run in its own goroutine.
func (c *Consumer) Start() {
	processed := 0 // Counter to keep track of the total number of tasks processed by this consumer.

	// The 'for range' loop over a channel will continuously receive values
	// until the channel is closed. Once the channel is closed (by the Pool)
	// and all values have been received, the loop will terminate gracefully.
	for task := range c.ResultChan {
		processed++ // Increment the counter for each task received.

		// Print the details of the processed task to the console.
		// This includes the task ID, its original data, and the calculated result (e.g., isPrime).
		fmt.Printf("Task   %d\t Data = %d\t isPrime = %v\n", task.ID, task.Data, task.Result)
	}

	// After the ResultChan is closed and all results have been consumed,
	// print a summary indicating the total number of tasks processed.
	fmt.Printf("Processed %d tasks\n", processed)
}
