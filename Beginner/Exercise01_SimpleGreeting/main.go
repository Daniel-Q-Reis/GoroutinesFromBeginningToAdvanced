package main

import (
	"fmt"
	"sync"
	"time"
)

// greet prints a greeting message for a specific name after a delay.
// It also signals the WaitGroup when done
func greet(name string, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the goroutine finishes

	fmt.Printf("Hello, my name is %s!\n", name)
	time.Sleep(time.Millisecond * 100) //Simulate some work
}

func main() {
	names := []string{"Alice", "Peter", "James", "Jordan", "Rob"}

	var wg sync.WaitGroup // Declare a WaitGroup to wait for goroutines to complete

	// Launch multiple goroutines
	for _, name := range names {
		wg.Add(1) // Increment the WaitGroup counter for each new goroutine
		go greet(name, &wg)
	}

	fmt.Println("Main goroutine continues...")

	wg.Wait() // Block until the WaitGroup counter is zero (all goroutines have called Done())

	fmt.Println("All greetings finished. Main goroutine exiting.")

}
