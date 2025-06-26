package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// cancellableWorker simulates a long-running goroutine that can be cancelled.
// It processes tasks from a channel or shuts down if the context is cancelled.
func cancelLableWorker(ctx context.Context, workerID int, taskChannel <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done(): // Case 1: Context cancellation signal received
			fmt.Printf("Worker %d: Context cancelled. Shutting down. \n", workerID)
			return //Exit the goroutine
		case task := <-taskChannel: // Case 2: A newtask is received
			fmt.Printf("Worker %d: Processing task %d\n", workerID, task)
			time.Sleep(time.Millisecond * 300)
		default: // Case 3: No task available and no cancellation signal
			// Important: prevents the goroutine from blocking indefinitely if no tasks or cancellation.
			// It allows the select to "poll" the context.Done() channel periodically.
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func main() {
	// Create a context with a cancel function.
	// 'ctx' is the context that workers will listen to.
	// 'cancel' is the function that will trigger the cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	numWorkers := 3
	taskChannel := make(chan int)
	var workerWg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		workerWg.Add(1)
		go cancelLableWorker(ctx, i, taskChannel, &workerWg)
	}

	// Goroutine to simulate sending tasks.
	go func() {
		for i := 0; i <= 10; i++ {
			taskChannel <- i
			time.Sleep(time.Millisecond * 150)
		}
		// Note: taskChannel is NOT closed here to demonstrate cancellation while tasks might still be in flight.
		// In a real scenario, you might close it if all tasks are sent and no more are expected
	}()

	// Goroutine to simulate an event that triggers cancellation.
	// The sleep duration ensures cancellation happens *during* task processing, in this case must be less than 1500 ms.
	go func() {
		time.Sleep(time.Millisecond * 1100)
		cancel()
		fmt.Println("Main: Context cancelled. Signaling workers to stop.")
	}()

	workerWg.Wait()

	fmt.Println("The program succed to stop workers using ''cancel() context.Context''")
}
