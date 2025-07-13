package exercise02workerpool

import (
	"fmt"
)

// Consumer processes results from the result channel
type Consumer struct {
	ResultChan <-chan Task
}

// NewConsumer creates a new consumer instance
func NewConsumer(resultChan <-chan Task) *Consumer {
	return &Consumer{
		ResultChan: resultChan,
	}
}

// Start begins consuming results
func (c *Consumer) Start() {

	processed := 0

	for task := range c.ResultChan {
		processed++
		fmt.Printf("Task %d: Data=%d, isPrime=%v\n", task.ID, task.Data, task.Result)
	}

	fmt.Printf("Processed %d tasks\n", processed)
}
