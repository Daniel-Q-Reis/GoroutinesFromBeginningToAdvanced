package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Task struct represents a task to be processed.
type Task struct {
	ID     int
	Number int
}

// worker processes tasks from the 'tasks' channel and sends results to the 'results' channel.
func worker(workerID int, tasks <-chan Task, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing Task ID %d (Number: %d)\n", workerID, task.ID, task.Number)
		result := task.Number * task.Number // Calculate square
		time.Sleep(time.Millisecond * 100)  //simulate work
		results <- result
	}
}

func main() {
	numTasks := 20
	numWorkers := 5
	tasks := make(chan Task)  // Channel for tasks
	results := make(chan int) //Channel for results
	var workerWg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		workerWg.Add(1)
		go worker(i, tasks, results, &workerWg)
	}

	// Go-routine to generate tasks and close the tasks channel
	go func() {
		for i := 1; i <= numTasks; i++ {
			random := rand.Intn(100) + 1
			task := Task{ID: i, Number: random}
			tasks <- task
		}
		close(tasks) //Close tasks channel when all tasks are sent
	}()

	// Go-routine to wait for all workers to finish and then close the results channel
	go func() {
		workerWg.Wait()
		close(results)
		fmt.Println("Results channel closed.")
	}()

	fmt.Println("Collecting results...")
	for result := range results {
		mathSqrt := math.Sqrt(float64(result)) //Unmake the square, to realise wich work was done
		fmt.Printf("Result: %d after processing the Number: %0.f\n", result, mathSqrt)
	}

	fmt.Println("All tasks are done.")

}
