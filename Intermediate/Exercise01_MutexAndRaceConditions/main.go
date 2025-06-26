package main

import (
	"fmt"
	"sync"
	// "time" // Not needed for this specific exercise, but useful for simulating work
)

var counter int   // Shared variable, prone to race conditions without protection
var mu sync.Mutex // Mutex to protect access to the 'counter' variable

// incrementCounter increments the global counter a fixed number of times.
// It uses a Mutex to ensure atomic operations on the shared 'counter'.
func incrementCounter(wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this goroutine has finished

	for i := 0; i < 1000; i++ {
		mu.Lock()   // Acquire the lock, blocking other goroutines from accessing 'counter'
		counter++   // Critical section: increment the shared counter
		mu.Unlock() // Release the lock, allowing other goroutines to acquire it
	}
	// fmt.Println("Worker finished incrementing.") // Can be uncommented for detailed logging, but will print 100 times
}

func main() {
	var wg sync.WaitGroup // WaitGroup to wait for all worker goroutines

	numWorkers := 100 // Number of goroutines that will increment the counter

	// Launch multiple worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Increment WaitGroup counter for each worker
		// Passing '&wg' by reference is crucial for workers to signal completion
		// 'i' is not used by incrementCounter, so no need for 'go func(num int){}(i)' here.
		go incrementCounter(&wg)
	}

	wg.Wait() // Block until all worker goroutines have finished

	// Print the final value of the counter.
	// With Mutex, this should consistently match numWorkers * 1000.
	fmt.Printf("Final Counter Value (expected %d) received: %d\n", numWorkers*1000, counter)
}
