package main

import (
	"fmt"
	"sync"
	"time"
)

// producer generates numbers and sends them to the channel.
// It closes the channel after sending all numbers to signal completion.
func producer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this producer goroutine has finished
	for num := 1; num <= 5; num++ {
		fmt.Printf("Producer %d sending: %d\n", id, num)
		time.Sleep(time.Millisecond * 50) // Simulate production time
		ch <- num
	}
}

// consumer reads numbers from the channel until it's closed.
// It signals the WaitGroup when it finishes consuming all values.
func consumer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()       // Signal that this consumer goroutine has finished
	for num := range ch { // Loop continues until the channel is closed and drained
		fmt.Printf("Consumer %d received: %d\n", id, num)
		time.Sleep(time.Millisecond * 150) // Simulate processing time
	}
	fmt.Printf("Consumer %d finished. \n", id) //Inform that consumer has completed its loop
}

func main() {
	ch := make(chan int)           // Unbuffered channel for communication
	var wgProducers sync.WaitGroup // WaitGroup to track producers
	var wgConsumers sync.WaitGroup // WaitGroup to track consumers

	numProducers := 2
	for i := 1; i <= numProducers; i++ {
		wgProducers.Add(1) //Increment counter for each producer
		go producer(i, ch, &wgProducers)
	}

	numConsumers := 2
	for i := 1; i <= numConsumers; i++ {
		wgConsumers.Add(1) //Increment counter for each producer
		go consumer(i, ch, &wgConsumers)
	}

	// Goroutine to close the channel AFTER all producers have finished.
	// This prevents "send on closed channel" panics.
	go func() {
		wgProducers.Wait() // Wait for all producers to complete their work
		close(ch)          // Safely close the channel
		fmt.Println("Channel closed by main orchestrator")
	}()

	wgConsumers.Wait() // Wait for all consumers to finish reading from the channel

	fmt.Println("All jobs are done!")
}
