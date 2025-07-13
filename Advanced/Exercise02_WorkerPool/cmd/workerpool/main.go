package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	exercise02workerpool "github.com/Daniel-Q-Reis/GoroutinesFromBeginningToAdvanced/Advanced/Exercise02_WorkerPool"
)

func main() {
	startTime := time.Now()

	// System configuration
	const numTasks = 1000
	numWorkers := runtime.NumCPU()

	// Create worker pool
	pool := exercise02workerpool.NewPool(numWorkers)

	var wg sync.WaitGroup

	// Start Workers
	pool.Start()

	// Create and start producer
	wg.Add(1)
	producer := exercise02workerpool.NewProducer(numTasks, pool.TaskChan)
	go func() {
		defer wg.Done()
		producer.Start()
	}()

	// Create and start consumer
	wg.Add(1)
	consumer := exercise02workerpool.NewConsumer(pool.ResultChan)
	go func() {
		defer wg.Done()
		consumer.Start()
	}()

	wg.Wait()

	// Display execution summary
	fmt.Printf("\nSystem Summary:\n")
	fmt.Printf("Total tasks processed: %d\n", numTasks)
	fmt.Printf("Numer of workers: %d\n", numWorkers)
	fmt.Printf("Total execution time: %v\n", time.Since(startTime))
}
