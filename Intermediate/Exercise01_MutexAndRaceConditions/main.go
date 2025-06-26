package main

import (
	"fmt"
	"sync"
)

var counter int
var mu sync.Mutex

func incrementCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
	fmt.Println("Worker finished incrementing")
}

func main() {
	var wg sync.WaitGroup

	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go incrementCounter(&wg)
	}

	wg.Wait()

	fmt.Printf("Final counter Value (expected %d) recived: %d\n", numWorkers*1000, counter)
}
