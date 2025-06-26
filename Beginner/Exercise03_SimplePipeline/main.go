package main

import "fmt"

// generateNumbers sends a sequence of numbers to the output channel.
// It closes the channel after sending all numbers.
func generateNumbers(out chan<- int) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("Generator sent: %d\n", i)
		out <- i
	}
	close(out) // Close the output channel to signal completion to the next stage
}

// duplicateNumbers reads numbers from the input channel, doubles them, and sends to the output channel.
// It closes its output channel when the input channel is closed and drained.
func duplicateNumbers(in <-chan int, out chan<- int) {
	for num := range in { // Loop continues until 'in' channel is closed and drained
		doubled := num * 2
		fmt.Printf("Duplicator received: %d, sending: %d\n", num, doubled)
		out <- doubled
	}
	close(out) // Close the output channel to signal completion to the next stage
}

// addFive reads numbers from the input channel, adds five, and sends to the output channel.
// It closes its output channel when the input channel is closed and drained.
func addFive(in <-chan int, out chan<- int) {
	for num := range in { // Loop continues until 'in' channel is closed and drained
		final := num + 5
		fmt.Printf("Adder received: %d, sending: %d\n", num, final)
		out <- final
	}
	close(out) // Close the output channel to signal completion to the next stage
}

func main() {
	// Create channels to connect the pipeline stages
	nums := make(chan int)         // Generator -> Duplicator
	doubledNums := make(chan int)  // Duplicator -> Adder
	finalResults := make(chan int) // Adder -> Main (final consumer)

	// Start pipeline stages as goroutines
	go generateNumbers(nums)
	go duplicateNumbers(nums, doubledNums)
	go addFive(doubledNums, finalResults)

	// Main goroutine acts as the final consumer, collecting and printing results.
	// The for-range loop will terminate when finalResults channel is closed by the adder.
	for result := range finalResults {
		fmt.Printf("Final Result: %d\n", result)
	}

	fmt.Println("All pipeline stages are completed")
}
