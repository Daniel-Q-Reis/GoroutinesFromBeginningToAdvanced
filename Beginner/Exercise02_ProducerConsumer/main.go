package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := 1; num <= 5; num++ {
		fmt.Printf("Producer %d sending: %d\n", id, num)
		time.Sleep(time.Millisecond * 50)
		ch <- num
	}
}

func consumer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Consumer %d received: %d\n", id, num)
		time.Sleep(time.Millisecond * 150)
	}
	fmt.Printf("Consumer %d finished. \n", id)
}

func main() {
	ch := make(chan int)
	var wgProducers sync.WaitGroup
	var wgConsumers sync.WaitGroup

	numProducers := 2
	for i := 1; i <= numProducers; i++ {
		wgProducers.Add(1)
		go producer(i, ch, &wgProducers)
	}

	numConsumers := 2
	for i := 1; i <= numConsumers; i++ {
		wgConsumers.Add(1)
		go consumer(i, ch, &wgConsumers)
	}

	go func() {
		wgProducers.Wait()
		close(ch)
		fmt.Println("Channel closed by main orchestrator")
	}()

	wgConsumers.Wait()

	fmt.Println("All jobs are done!")
}
