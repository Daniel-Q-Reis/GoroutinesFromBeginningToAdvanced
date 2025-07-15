package main // The 'main' package indicates this is an executable program.

import (
	"fmt"     // Package for formatted I/O, used for printing output to the console.
	"runtime" // Provides functions to interact with the Go runtime, e.g., NumCPU.
	"sync"    // Package for synchronization primitives, e.g., WaitGroup.
	"time"    // Package for time-related functions, used for measuring execution time.

	// Import the 'exercise02workerpool' package, which contains all the core logic
	// for the worker pool components (Task, Worker, Producer, Consumer, Pool).
	// The import path reflects the module path and directory structure.
	exercise02workerpool "github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced/Advanced/Exercise02_WorkerPool"
)

// main is the entry point of the application.
func main() {
	// Record the start time to measure the total execution duration of the program.
	startTime := time.Now()

	// --- System Configuration ---
	const numTasks = 1000 // Define the total number of tasks to be generated and processed.
	// Determine the number of workers based on the number of available CPU cores.
	// This is a common practice to optimize CPU-bound workloads, allowing one worker
	// per core to maximize parallel execution without excessive context switching overhead.
	numWorkers := runtime.NumCPU()

	// --- Worker Pool Setup ---
	// Create a new instance of the worker Pool.
	// The pool will manage the workers and the task/result channels.
	pool := exercise02workerpool.NewPool(numWorkers)

	// A WaitGroup for the main function to synchronize the completion of the Producer
	// and Consumer goroutines. This is distinct from the internal WaitGroup used by the Pool.
	var wg sync.WaitGroup

	// --- Start Workers ---
	// Start the worker pool. This method launches 'numWorkers' goroutines,
	// each running a Worker.Start() loop, and also a goroutine that waits for
	// all workers to finish before closing the ResultChan.
	pool.Start()

	// --- Start Producer ---
	// Increment the main WaitGroup counter as the Producer will run in a separate goroutine.
	wg.Add(1)
	// Create a new Producer instance. It is given the total number of tasks to generate
	// and the TaskChan from the pool to send tasks to.
	producer := exercise02workerpool.NewProducer(numTasks, pool.TaskChan)
	// Launch the producer's Start method in a new goroutine.
	go func() {
		// Defer wg.Done() ensures the main WaitGroup counter is decremented when
		// the producer goroutine finishes its execution (after sending all tasks and closing TaskChan).
		defer wg.Done()
		producer.Start() // The producer starts generating and sending tasks.
	}()

	// --- Start Consumer ---
	// Increment the main WaitGroup counter for the Consumer goroutine.
	wg.Add(1)
	// Create a new Consumer instance. It is given the ResultChan from the pool
	// to receive processed tasks from.
	consumer := exercise02workerpool.NewConsumer(pool.ResultChan)
	// Launch the consumer's Start method in a new goroutine.
	go func() {
		// Defer wg.Done() ensures the main WaitGroup counter is decremented when
		// the consumer goroutine finishes (after the ResultChan is closed and drained).
		defer wg.Done()
		consumer.Start() // The consumer starts receiving and displaying results.
	}()

	// --- Await Completion ---
	// wg.Wait() blocks the main goroutine until all goroutines associated with
	// the main WaitGroup (i.e., Producer and Consumer) have called wg.Done().
	// This ensures that all tasks are generated, processed, and consumed before
	// the main function proceeds to display the summary or exits.
	wg.Wait()

	// --- Display Execution Summary ---
	// Calculate the total time elapsed since the program started.
	elapsedTime := time.Since(startTime)
	fmt.Printf("\nSystem Summary:\n")
	fmt.Printf("Total tasks processed: %d\n", numTasks)   // This refers to the number of tasks the producer was configured to generate.
	fmt.Printf("Number of workers: %d\n", numWorkers)     // Displays the number of workers utilized.
	fmt.Printf("Total execution time: %v\n", elapsedTime) // Displays the total time taken for the entire process.
}
