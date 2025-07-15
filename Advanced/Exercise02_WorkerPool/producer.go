package exercise02workerpool

import (
	"math/rand" // Package for generating pseudo-random numbers.
	"time"      // Package for time-related functions, used for seeding the random number generator.
)

// Producer is responsible for generating tasks and sending them to the task channel.
type Producer struct {
	TaskCount    int         // The total number of tasks this producer will generate.
	TaskChan     chan<- Task // A send-only channel where the producer sends newly created tasks.
	RandomNumber *rand.Rand  // A source of pseudo-random numbers for generating task data and complexity.
}

// NewProducer creates and returns a new Producer instance.
// It initializes the producer with the total number of tasks to create
// and the channel through which it will send these tasks.
func NewProducer(taskCount int, taskChan chan<- Task) *Producer {
	return &Producer{
		TaskCount: taskCount, // Sets the total number of tasks to be generated.
		TaskChan:  taskChan,  // Assigns the channel to send tasks.
		// Initializes a new pseudo-random number generator.
		// rand.NewSource(time.Now().UnixNano()) seeds the generator with the current nanosecond timestamp,
		// ensuring different sequences of random numbers on each program run.
		RandomNumber: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins the task generation process.
// This method is designed to be run in its own goroutine.
func (p *Producer) Start() {
	// Loop 'TaskCount' times to generate the specified number of tasks.
	for i := 0; i < p.TaskCount; i++ {
		// Create a new Task instance for each iteration.
		task := Task{
			ID: i, // Assigns a sequential ID to the task.

			// Generates a random integer for the task's data.
			// Intn(1000000) generates numbers from 0 to 999999. Adding 1 makes it 1 to 1000000.
			Data: p.RandomNumber.Intn(1000000) + 1,

			// Generates a random complexity (simulated processing time) for the task.
			// Intn(195) generates numbers from 0 to 194. Adding 5 makes it 5 to 199.
			// Multiplied by time.Millisecond to convert to a time.Duration.
			Complexity: time.Duration(p.RandomNumber.Intn(195)+5) * time.Millisecond,
		}

		// Send the newly created task to the TaskChan.
		// Since TaskChan is unbuffered (as defined in Pool), this send operation
		// will block if no worker is ready to receive the task. This mechanism
		// provides backpressure, preventing the producer from overwhelming the workers.
		p.TaskChan <- task
	}

	// After all tasks have been sent, close the TaskChan.
	// Closing the channel signals to all listening workers (via 'for range' loops)
	// that no more tasks will be sent, allowing them to gracefully exit their loops.
	close(p.TaskChan)
}
