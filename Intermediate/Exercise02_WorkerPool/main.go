package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID     int
	Number int
}

func worker(workerID int, tasks <-chan Task, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing Task ID %d (Number: %d)\n", workerID, task.ID, task.Number)
		result := task.Number * task.Number
		time.Sleep(time.Millisecond * 100)
		results <- result
	}
}

func main() {
	numTasks := 20
	numWorkers := 5
	tasks := make(chan Task)
	results := make(chan int)
	var workerWg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		workerWg.Add(1)
		go worker(i, tasks, results, &workerWg)
	}

	go func() {
		for i := 1; i <= numTasks; i++ {
			random := rand.Intn(100) + 1
			task := Task{ID: i, Number: random}
			tasks <- task
		}
		close(tasks)
	}()

	go func() {
		workerWg.Wait()
		close(results)
		fmt.Println("Results channel closed.")
	}()

	fmt.Println("Collecting results...")
	for result := range results {
		mathSqrt := math.Sqrt(float64(result))
		fmt.Printf("Result: %d after processing the Number: %0.f\n", result, mathSqrt)
	}

	fmt.Println("All tasks are done.")

}
