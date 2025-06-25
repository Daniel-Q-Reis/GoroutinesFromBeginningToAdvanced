package main

import "fmt"

func generateNumbers(out chan<- int) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("Generator sent: %d\n", i)
		out <- i
	}
	close(out)
}

func duplicateNumbers(in <-chan int, out chan<- int) {
	for num := range in {
		doubled := num * 2
		fmt.Printf("Duplicator received: %d, sending: %d\n", num, doubled)
		out <- doubled
	}
	close(out)
}

func addFive(in <-chan int, out chan<- int) {
	for num := range in {
		final := num + 5
		fmt.Printf("Adder received: %d, sending: %d\n", num, final)
		out <- final
	}
	close(out)
}

func main() {
	nums := make(chan int)
	doubledNums := make(chan int)
	finalResults := make(chan int)

	go generateNumbers(nums)
	go duplicateNumbers(nums, doubledNums)
	go addFive(doubledNums, finalResults)

	for result := range finalResults {
		fmt.Printf("Final Result: %d\n", result)
	}

	fmt.Println("All pipeline stages are completed")
}
